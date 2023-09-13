package main

import (
	"github.com/mangenotwork/beacon-tower/udp"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
	"time"
)

func main() {
	conf.InitConf("./conf/")
	master, err := conf.YamlGetString("master")
	if err != nil {
		panic(err)
	}
	log.Info(master)
	client, err := udp.NewClient(master)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(2 * time.Second)
			rse, err := client.Get("conn/test", []byte("test"))
			if err != nil {
				udp.Error(err)
				return
			}
			udp.Info("get 请求返回 = ", string(rse))
		}
	}()

	// 运行客户端
	client.Run()
}
