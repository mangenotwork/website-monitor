package main

import (
	"net/http"
	"time"

	"website-monitor/master/business"
	"website-monitor/master/dao"
	"website-monitor/master/routers"
	"website-monitor/master/udpServer"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
)

func main() {
	conf.InitConf("./conf/")
	dao.DB.Init()
	// 启动 udp servers
	go udpServer.RunUDPServer()
	// 启动采集
	go business.Collect()
	// 启动 https servers
	gin.SetMode(gin.DebugMode)
	server := &http.Server{
		Addr:           ":" + conf.Conf.Default.HttpServer.Prod,
		Handler:        routers.Routers(),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
