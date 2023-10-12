package handler

import (
	"github.com/mangenotwork/beacon-tower/udp"
	"github.com/mangenotwork/common/log"
	"website-monitor/monitor/business"
)

func NoticeUpdateWebsite(c *udp.Client, data []byte) {
	log.Info("更新监测网站")
	log.Info("获取参数: ", string(data))
	business.GetWebsite()
}

func NoticeUpdateWebsiteAllUrl(c *udp.Client, data []byte) {
	log.Info("更新监测网站url")
	log.Info("获取参数: ", string(data))
}

func NoticeUpdateWebsitePoint(c *udp.Client, data []byte) {
	log.Info("更新监测网站监测点")
	log.Info("获取参数: ", string(data))
	hostID := string(data)
	business.GetWebsitePoint(hostID)
}
