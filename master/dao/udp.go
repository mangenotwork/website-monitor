package dao

import (
	"github.com/mangenotwork/beacon-tower/udp"
	"website-monitor/master/constname"
)

var Servers *udp.Servers

// NoticeUpdateWebsite 有新的监测网站，通知拉取并更新监测网站列表
func NoticeUpdateWebsite() {
	Servers.NoticeAll(constname.NoticeUpdateWebsiteLabel, []byte(""), Servers.SetNoticeRetry(2, 3000))
}

// NoticeUpdateWebsitePoint 有网站监测点更改，通知拉取并更新该网站监测点
func NoticeUpdateWebsitePoint(hostID string) {
	Servers.NoticeAll(constname.NoticeUpdateWebsitePointLabel, []byte(hostID), Servers.SetNoticeRetry(2, 3000))
}

// TODO... 有网站监测设置更改，通知拉取并更新该网站的监测设置

// TODO... 通知执行压力并发请求任务
