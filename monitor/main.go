package main

import (
	"github.com/mangenotwork/beacon-tower/udp"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
	"website-monitor/monitor/business"
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
	// 初始化
	business.Initialize(client)
	// 业务
	business.Business(client)
	// 运行客户端
	client.Run()
}
