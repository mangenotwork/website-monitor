package udpServer

import (
	"fmt"

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

	// 定义get方法
	Servers.GetHandleFunc("conn/test", ConnTest)

	// 启动servers
	Servers.Run()
}

func ConnTest(s *udp.Servers, param []byte) (int, []byte) {
	udp.Info("获取到的请求参数  param = ", string(param))
	return 0, []byte(fmt.Sprintf("服务器名称 %s.", s.GetServersName()))
}
