package dao

import (
	"github.com/mangenotwork/beacon-tower/udp"
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
	Key         string `json:"key"`
	Name        string `json:"name"`
	Online      bool   `json:"online"`
	Addr        string `json:"addr"`
	LastTime    int64  `json:"lastTime"`
	DiscardTime int64  `json:"discardTime"`
	// TODO 获取ip环境地址
	// TODO 获取系统信息
}

// GetClientList 获取监测器列表信息
func GetClientList() []*MonitorInfo {
	data := make([]*MonitorInfo, 0)
	table := Servers.OnLineTable()
	for k, v := range table {
		data = append(data, &MonitorInfo{
			Key:         k,
			Name:        v.Name,
			Online:      v.Online,
			Addr:        v.Addr,
			LastTime:    v.LastTime,
			DiscardTime: v.DiscardTime,
		})
	}
	return data
}
