package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
	"net/http"
	"time"
	"website-monitor/master/routers"
	"website-monitor/master/udpServer"
)

func main() {
	conf.InitConf("./conf/")

	// 启动 udp servers
	go udpServer.RunUDPServer()

	// 启动 https servers
	gin.SetMode(gin.ReleaseMode)
	server := &http.Server{
		Addr:           ":" + conf.Conf.Default.HttpServer.Prod,
		Handler:        routers.Routers(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
