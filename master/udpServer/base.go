package udpServer

import (
	"github.com/mangenotwork/beacon-tower/udp"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/utils"
)

var Servers *udp.Servers

func RunUDPServer() {
	var err error
	// 初始化 s端
	Servers, err = udp.NewServers("0.0.0.0", utils.AnyToInt(conf.Conf.Default.UdpServer.Prod))
	if err != nil {
		panic(err)
	}

	// 定义put方法
	Servers.PutHandleFunc("monitor", MonitorRse)

	// 启动servers
	Servers.Run()
}

func MonitorRse(s *udp.Servers, param []byte) {
	udp.Info("接收到监测结果  param = ", string(param))
}
