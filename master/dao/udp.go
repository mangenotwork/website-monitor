package dao

import (
	"encoding/json"
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

// GetClientList 获取监测器列表信息
func GetClientList() []*MonitorInfo {
	data := make([]*MonitorInfo, 0)
	table := Servers.OnLineTable()

	for k, v := range table {

		info := &MonitorInfo{
			Key:         k,
			Name:        v.Name,
			Online:      v.Online,
			IP:          v.IP,
			Addr:        v.Addr,
			LastTime:    v.LastTime,
			DiscardTime: v.DiscardTime,
		}

		info.PublicIP, info.PublicIPAddr = getIPAddr(v.Name, v.IP)
		info.OSInfo = getOSInfo(v.Name, v.IP)
		data = append(data, info)
	}

	return data
}

func getIPAddr(name, ip string) (string, string) {
	data, err := Servers.GetAtIP("ipAddr", name, ip, []byte(""))
	if err != nil {
		log.Error(err)
	}

	ipInfo := &IPInfo{}
	err = json.Unmarshal(data, &ipInfo)
	if err != nil {
		log.Error(err)
	}

	log.Info("ipInfo = ", ipInfo)
	return ipInfo.IP, ipInfo.Address
}

type OSInfo struct {
	HostName      string `json:"hostName"`
	OSType        string `json:"osType"`
	OSArch        string `json:"osArch"`
	CpuCoreNumber string `json:"cpuCoreNumber"`
	InterfaceInfo string `json:"interfaceInfo"`
}

func getOSInfo(name, ip string) *OSInfo {
	data, err := Servers.GetAtIP("osInfo", name, ip, []byte(""))
	if err != nil {
		log.Error(err)
	}

	osInfo := &OSInfo{}
	err = json.Unmarshal(data, &osInfo)
	if err != nil {
		log.Error(err)
	}

	log.Info("osInfo = ", osInfo)
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
