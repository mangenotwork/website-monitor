package handler

import (
	"github.com/mangenotwork/beacon-tower/udp"
	"github.com/mangenotwork/common/log"
	"website-monitor/monitor/business"
)

func NoticeUpdateWebsiteAll(c *udp.Client, data []byte) {
	log.Info("更新监测网站")
	log.Info("获取参数: ", string(data))
	business.GetWebsiteAll()
}

func NoticeUpdateWebsite(c *udp.Client, data []byte) {
	log.Info("更新监测网站指定")
	log.Info("获取参数: ", string(data))
	hostID := string(data)
	business.GetWebsite(hostID)
}

func NoticeDelWebsite(c *udp.Client, data []byte) {
	log.Info("删除监测网站指定")
	log.Info("获取参数: ", string(data))
	hostID := string(data)
	business.DelWebsite(hostID)
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
