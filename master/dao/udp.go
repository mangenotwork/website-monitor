package dao

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/mangenotwork/common/log"
	udp "github.com/mangenotwork/udp_comm"
	"website-monitor/master/constname"
)

var Servers *udp.Servers

// NoticeUpdateWebsite 有新的监测网站，通知拉取并更新监测网站列表
func NoticeUpdateWebsite(hostID string) {
	Servers.NoticeAll(constname.NoticeUpdateWebsiteLabel, []byte(hostID), Servers.SetNoticeRetry(2, 3000))
}

func NoticeDelWebsite(hostID string) {
	Servers.NoticeAll(constname.NoticeDelWebsiteLabel, []byte(hostID), Servers.SetNoticeRetry(2, 3000))
}

// NoticeUpdateWebsitePoint 有网站监测点更改，通知拉取并更新该网站监测点
func NoticeUpdateWebsitePoint(hostID string) {
	Servers.NoticeAll(constname.NoticeUpdateWebsitePointLabel, []byte(hostID), Servers.SetNoticeRetry(2, 3000))
}

// TODO... 有网站监测设置更改，通知拉取并更新该网站的监测设置

// TODO... 通知执行压力并发请求任务

type MonitorInfo struct {
	Key          string  `json:"key"`
	Name         string  `json:"name"`
	Online       bool    `json:"online"`
	IP           string  `json:"IP"`
	Addr         string  `json:"addr"`
	LastTime     int64   `json:"lastTime"`
	DiscardTime  int64   `json:"discardTime"`
	PublicIP     string  `json:"publicIP"`     // 公网ip环境地址
	PublicIPAddr string  `json:"publicIPAddr"` // 公网ip属地
	OSInfo       *OSInfo `json:"osInfo"`       // TODO 获取系统信息
}

// HasOnlineMonitor 只读内存中的在线表，不发起任何 UDP 请求。
// 用于「是否有在线探针」这类判断，避免为了一个布尔值去同步拉取每个探针的 IP/OS。
func HasOnlineMonitor() bool {
	if Servers == nil {
		return false
	}
	for _, v := range Servers.OnLineTable() {
		if v.Online {
			return true
		}
	}
	return false
}

// GetClientList 轻量列表：只读内存中的在线表，不发起任何 UDP 请求。
// 首页侧栏这类只需要名称/在线状态/IP 的场景用这个，避免被探针的 1s 超时拖住。
func GetClientList() []*MonitorInfo {
	if Servers == nil {
		return make([]*MonitorInfo, 0)
	}
	table := Servers.OnLineTable()
	data := make([]*MonitorInfo, 0, len(table))

	for k, v := range table {
		data = append(data, &MonitorInfo{
			Key:         k,
			Name:        v.Name,
			Online:      v.Online,
			IP:          v.IP,
			Addr:        v.Addr,
			LastTime:    v.LastTime,
			DiscardTime: v.DiscardTime,
		})
	}
	return data
}

// 探针详情缓存：IP/系统信息变化不频繁，避免每次打开探针页都串行等 UDP
type monitorEnrich struct {
	publicIP   string
	publicAddr string
	osInfo     *OSInfo
	at         time.Time
}

var monitorEnrichCache sync.Map

const (
	monitorEnrichTTL     = 5 * time.Minute
	monitorEnrichFailTTL = 30 * time.Second // 失联节点的空结果短缓存，避免每次刷新都等超时
	monitorEnrichBudget  = 2 * time.Second  // 单次请求最多花在 UDP 上的时间
	monitorEnrichTimeout = 800              // 单个 UDP Get 上限，单位 ms
)

// GetClientListDetail 需要公网 IP / 系统信息时使用，只对在线节点发起 UDP 请求。
//
// 注意：udp_comm v0.0.1 的 Get 请求 ID 由每次新建的 IdWorker 生成，并发调用会
// 出现 ID 冲突（实测 4 个并发请求 3 个超时、响应错配），因此这里必须串行。
// 为避免失联节点把页面拖住，这里加了 5 分钟缓存 + 单次 2s 总预算。
func GetClientListDetail() []*MonitorInfo {
	data := GetClientList()
	deadline := time.Now().Add(monitorEnrichBudget)

	for _, info := range data {
		if !info.Online {
			continue
		}
		key := info.Name + "|" + info.IP

		if v, ok := monitorEnrichCache.Load(key); ok {
			cached := v.(monitorEnrich)
			info.PublicIP = cached.publicIP
			info.PublicIPAddr = cached.publicAddr
			info.OSInfo = cached.osInfo
			ttl := monitorEnrichTTL
			if cached.publicIP == "" && cached.osInfo == nil {
				ttl = monitorEnrichFailTTL
			}
			if time.Since(cached.at) < ttl {
				continue
			}
		}

		remaining := int(time.Until(deadline) / time.Millisecond)
		if remaining < 100 {
			continue // 预算耗尽，保留缓存值（可能为空），下次刷新再补
		}
		timeout := remaining
		if timeout > monitorEnrichTimeout {
			timeout = monitorEnrichTimeout
		}
		half := timeout / 2

		ip, addr := getIPAddr(info.Name, info.IP, half)
		osInfo := getOSInfo(info.Name, info.IP, timeout-half)
		if ip == "" && osInfo == nil {
			// 失败也记一个短缓存，30s 内不再为同一节点重复等超时
			monitorEnrichCache.Store(key, monitorEnrich{at: time.Now()})
			continue
		}

		info.PublicIP, info.PublicIPAddr, info.OSInfo = ip, addr, osInfo
		monitorEnrichCache.Store(key, monitorEnrich{
			publicIP: ip, publicAddr: addr, osInfo: osInfo, at: time.Now(),
		})
	}
	return data
}

func getIPAddr(name, ip string, timeout int) (string, string) {
	data, err := Servers.GetAtIPTimeOut(timeout, "ipAddr", name, ip, []byte(""))
	if err != nil {
		log.Debug("获取探针 ipAddr 失败: ", err)
		return "", ""
	}

	ipInfo := &IPInfo{}
	if err := json.Unmarshal(data, ipInfo); err != nil {
		log.Debug("解析探针 ipAddr 失败: ", err)
		return "", ""
	}
	return ipInfo.IP, ipInfo.Address
}

type OSInfo struct {
	HostName      string `json:"hostName"`
	OSType        string `json:"osType"`
	OSArch        string `json:"osArch"`
	CpuCoreNumber string `json:"cpuCoreNumber"`
	InterfaceInfo string `json:"interfaceInfo"`
}

func getOSInfo(name, ip string, timeout int) *OSInfo {
	data, err := Servers.GetAtIPTimeOut(timeout, "osInfo", name, ip, []byte(""))
	if err != nil {
		log.Debug("获取探针 osInfo 失败: ", err)
		return nil
	}

	osInfo := &OSInfo{}
	if err := json.Unmarshal(data, osInfo); err != nil {
		log.Debug("解析探针 osInfo 失败: ", err)
		return nil
	}
	return osInfo
}

func GetClientList2() {
	table := Servers.OnLineTable()

	for k, v := range table {
		log.Info(k, v)
		data, err := Servers.GetAtIP("ipAddr", v.Name, v.IP, []byte(""))
		if err != nil {
			log.Error(err)
		}

		inInfo := &IPInfo{}
		err = json.Unmarshal(data, &inInfo)
		if err != nil {
			log.Error(err)
		}

		log.Info("inInfo = ", inInfo)
	}

}
