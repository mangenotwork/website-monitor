package udpServer

import (
	"github.com/mangenotwork/beacon-tower/udp"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	"io"
	"os"
	"strings"
	"website-monitor/master/entity"
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

func MonitorRse(s *udp.Servers, c *udp.ClientInfo, param []byte) {
	udp.Info("接收到监测结果  param = ", string(param))
	ip := c.Addr.String()
	mLogStr := string(param) + ip + "|"
	mLog := toMonitorLogObj(ip, mLogStr)
	// 写日志
	WriteLog(mLog.HostId, mLogStr)
	// TODO 报警统计
}

const globalDayLayout = "20060102"

func toMonitorLogObj(ip, str string) *entity.MonitorLog {
	strList := strings.Split(str, "|")
	if len(strList) < 17 {
		return nil
	}
	return &entity.MonitorLog{
		LogType:         strList[0],
		Time:            strList[1],
		HostId:          strList[2],
		Host:            strList[3],
		UriType:         strList[4],
		Uri:             strList[5],
		UriCode:         utils.AnyToInt(strList[6]),
		UriMs:           utils.AnyToInt64(strList[7]),
		ContrastUri:     strList[8],
		ContrastUriCode: utils.AnyToInt(strList[9]),
		ContrastUriMs:   utils.AnyToInt64(strList[10]),
		Ping:            strList[11],
		PingMs:          utils.AnyToInt64(strList[12]),
		Msg:             strList[13],
		MonitorName:     strList[14],
		MonitorIP:       strList[15],
		MonitorAddr:     strList[16],
		ClientIP:        ip,
	}
}

// WriteLog 写日志
func WriteLog(hostId, mLog string) {
	logPath, err := conf.YamlGetString("logPath")
	if err != nil {
		logPath = "./log/"
	}
	fileName := logPath + hostId + "_" + utils.NowDateLayout(globalDayLayout) + ".log"
	//log.Info("fileName = ", fileName)
	var file *os.File
	if !utils.Exists(fileName) {
		_ = os.MkdirAll(logPath, 0666)
		file, _ = os.Create(fileName)
	} else {
		file, _ = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = io.WriteString(file, mLog+"\n")
	if err != nil {
		log.Error("写入日志错误：", err)
		return
	}
}
