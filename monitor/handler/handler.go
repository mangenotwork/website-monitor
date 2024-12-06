package handler

import (
	"encoding/json"
	"github.com/mangenotwork/common/log"
	udp "github.com/mangenotwork/udp_comm"
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

func GetIPAddr(c *udp.Client, param []byte) (int, []byte) {
	log.Info("获取ip地址与属地")
	data := business.GetMyIP()
	log.Info("data = ", data)

	b, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
	}

	return 0, b
}

func GetOSInfo(c *udp.Client, param []byte) (int, []byte) {
	log.Info("获取监测器宿主系统信息")
	data := &OSInfo{
		HostName:      GetHostName(),
		OSType:        GetSysType(),
		OSArch:        GetSysArch(),
		CpuCoreNumber: GetCpuCoreNumber(),
		InterfaceInfo: GetInterfaceInfo(),
	}

	b, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
	}

	return 0, b
}
