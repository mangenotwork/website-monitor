package main

import (
	"github.com/mangenotwork/beacon-tower/udp"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
	"website-monitor/monitor/business"
	"website-monitor/monitor/handler"
)

func main() {
	conf.InitConf("./conf/")
	master, err := conf.YamlGetString("master")
	if err != nil {
		panic(err)
	}
	business.MasterHTTP, err = conf.YamlGetString("masterHTTP")
	if err != nil {
		panic(err)
	}
	log.Info(master)
	client, err := udp.NewClient(master)
	if err != nil {
		panic(err)
	}

	// 初始化监测
	business.Initialize(client)

	// 通知方法
	client.NoticeHandleFunc("website", handler.NoticeUpdateWebsite)      // 通知更新网站监测
	client.NoticeHandleFunc("allUrl", handler.NoticeUpdateWebsiteAllUrl) // 通知更新网站url
	client.NoticeHandleFunc("point", handler.NoticeUpdateWebsitePoint)   // 通知更新网站监测点
	// TODO 通知执行并发请求

	// 运行客户端
	client.Run()
}
