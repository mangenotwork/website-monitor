package udpServer

import (
	"github.com/mangenotwork/beacon-tower/udp"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/utils"
	"website-monitor/master/dao"
)

func RunUDPServer() {
	var err error
	connCode, err := conf.YamlGetString("connCode")
	if err != nil {
		connCode = udp.DefaultConnectCode
	}
	connSecret, err := conf.YamlGetString("connSecret")
	if err != nil {
		connSecret = udp.DefaultSecretKey
	}
	// 初始化 s端
	dao.Servers, err = udp.NewServers("0.0.0.0",
		utils.AnyToInt(conf.Conf.Default.UdpServer.Prod),
		udp.SetServersConf("s", connCode, connSecret))
	if err != nil {
		panic(err)
	}

	// 定义put方法
	dao.Servers.PutHandleFunc("monitor", MonitorRse) // 监测结果
	dao.Servers.PutHandleFunc("stress", StressRse)   // TODO... 并发请求结果

	// 启动servers
	dao.Servers.Run()
}

func MonitorRse(s *udp.Servers, c *udp.ClientInfo, param []byte) {
	udp.Info("接收到监测结果  param = ", string(param))
	ip := c.Addr.String()
	mLogStr := string(param) + ip + "|"
	mLogDao := dao.NewMonitorLogDao()
	mLog := mLogDao.ToMonitorLogObj(mLogStr)
	// 写日志
	mLogDao.Write(mLog.HostId, mLogStr)
	// TODO 报警统计
}

func StressRse(s *udp.Servers, c *udp.ClientInfo, param []byte) {

}
