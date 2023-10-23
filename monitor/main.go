package main

import (
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
	udp "github.com/mangenotwork/udp_comm"
	"math/rand"
	"time"
	"website-monitor/monitor/business"
	"website-monitor/monitor/handler"
)

const maxNameLen = 7

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	clientName, err := conf.YamlGetString("clientName")
	if err != nil {
		clientName = randStringBytes(maxNameLen)
	}
	connectCode, err := conf.YamlGetString("connCode")
	if err != nil {
		connectCode = udp.DefaultConnectCode
	}
	secretKey, err := conf.YamlGetString("connSecret")
	if err != nil {
		secretKey = udp.DefaultSecretKey
	}
	client, err := udp.NewClient(master, udp.ClientConf{
		Name:        clientName,
		ConnectCode: connectCode,
		SecretKey:   secretKey,
	})
	if err != nil {
		panic(err)
	}

	// 初始化监测
	business.Initialize(client)

	// 获取方法
	client.GetHandleFunc("ipAddr", handler.GetIPAddr) // 获取监测器ip属地
	client.GetHandleFunc("osInfo", handler.GetOSInfo) // 获取监测器宿主系统信息

	// 通知方法
	client.NoticeHandleFunc("websiteAll", handler.NoticeUpdateWebsiteAll) // 通知更新网站监测-所有
	client.NoticeHandleFunc("website", handler.NoticeUpdateWebsite)       // 通知更新网站监测-指定
	client.NoticeHandleFunc("websiteDel", handler.NoticeDelWebsite)       // 通知删除网站监测-指定
	client.NoticeHandleFunc("allUrl", handler.NoticeUpdateWebsiteAllUrl)  // 通知更新网站url
	client.NoticeHandleFunc("point", handler.NoticeUpdateWebsitePoint)    // 通知更新网站监测点
	// TODO 通知执行并发请求

	// 运行客户端
	client.Run()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
