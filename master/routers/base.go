package routers

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/ginHelper"
	"net/http"
	"website-monitor/master/handler"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
}

func Routers() *gin.Engine {
	Router.Use(gzip.Gzip(gzip.DefaultCompression))
	Router.StaticFS("/static", http.Dir("./static"))
	Router.Delims("{[", "]}")
	//Router.LoadHTMLGlob("views/**/*")
	Router.GET("/", ginHelper.Handle(handler.Index))
	return Router
}
