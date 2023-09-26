package routers

import (
	"net/http"

	"website-monitor/master/constname"
	"website-monitor/master/handler"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/ginHelper"
	"github.com/mangenotwork/common/utils"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
}

func Routers() *gin.Engine {
	Router.Use(gzip.Gzip(gzip.DefaultCompression))
	Router.StaticFS("/static", http.Dir("./static"))
	Router.Delims("{[", "]}")
	Svg()
	Router.LoadHTMLGlob("views/**/*")
	Login()
	Page()
	API()
	Data()
	Test()
	return Router
}

func Login() {
	login := Router.Group("")
	login.Use(ginHelper.CSRFMiddleware())
	login.GET("/", handler.LoginPage)
	login.POST("/login", handler.Login)
}

func Page() {
	// 404 && 405 && err page
	Router.NoRoute(handler.NotFond)
	Router.NoMethod(handler.NotFond)
	// page group
	pg := Router.Group("")
	pg.Use(AuthPG())
	pg.GET("/home", handler.HomePage)
	pg.GET("/monitor", handler.MonitorPage)
	pg.GET("/tool", handler.ToolPage)
	pg.GET("/test/api", handler.TestAPI)
	pg.GET("/test/stress", handler.TestStress)
	pg.GET("/test/penetration", handler.TestPenetration)
	pg.GET("/operation", handler.Operation)
	pg.GET("/instructions", handler.Instructions)
}

func API() {
	api := Router.Group("/api").Use(AuthAPI())
	api.GET("/out", handler.Out)

	// website
	api.POST("/website/add", ginHelper.Handle(handler.WebsiteAdd))                 // 创建网站监测
	api.GET("/website/list", ginHelper.Handle(handler.WebsiteList))                // 监测网站列表
	api.GET("/website/info/:host", ginHelper.Handle(handler.WebsiteInfo))          // 监测网站详情
	api.GET("/website/info/refresh", ginHelper.Handle(handler.WebsiteInfoRefresh)) // 监测网站详情刷新
	api.GET("/website/delete/:host", ginHelper.Handle(handler.WebsiteDelete))      // 删除网站监测
	api.GET("/website/urls/:host", ginHelper.Handle(handler.WebsiteUrls))          // 监测网站采集到url
	api.POST("/website/edit", ginHelper.Handle(handler.WebsiteEdit))               // TODO 监测设置
	api.GET("/website/chart/:host", ginHelper.Handle(handler.WebsiteChart))        // TODO 图表
	api.GET("/website/alert/:host", ginHelper.Handle(handler.WebsiteAlertList))    // TODO 报警信息
	api.GET("/website/alert/del/:host", ginHelper.Handle(handler.WebsiteAlertDel)) // TODO 报警信息

	// mail
	api.GET("/mail/init", ginHelper.Handle(handler.MailInit))          // 是否设置邮件
	api.POST("/mail/conf", ginHelper.Handle(handler.MailConf))         // 设置邮件配置
	api.GET("/mail/info", ginHelper.Handle(handler.MailInfo))          // 获取邮件配置信息
	api.POST("/mail/sendTest", ginHelper.Handle(handler.MailSendTest)) // 测试发生邮件

	// tool
	api.POST("/tool/history", ginHelper.Handle(handler.ToolHistorySet))            // 记录历史记录
	api.GET("/tool/history", ginHelper.Handle(handler.ToolHistoryGet))             // 获取历史记录
	api.GET("/tool/history/clear", ginHelper.Handle(handler.ToolHistoryClear))     // 清空历史记录
	api.GET("/tool/certificate", ginHelper.Handle(handler.GetSSLCertificate))      // 获取证书
	api.GET("/tool/nsLookUp", ginHelper.Handle(handler.DNSLookUp))                 // 查询dns
	api.GET("/tool/nsLookUp/all", ginHelper.Handle(handler.DNSLookUpAll))          // 查询dns
	api.GET("/tool/whois", ginHelper.Handle(handler.Whois))                        // Whois查询
	api.GET("/tool/ip", ginHelper.Handle(handler.IPInfo))                          // ip信息查询
	api.GET("/tool/myIP", ginHelper.Handle(handler.MyIPInfo))                      // 本机ip信息
	api.GET("/tool/website/tdki", ginHelper.Handle(handler.GetWebSiteTDKI))        // 获取网站的T, D, K, 图标
	api.GET("/tool/website/collectInfo", ginHelper.Handle(handler.CollectWebSite)) // 采集网站信息
	api.GET("/tool/icp", ginHelper.Handle(handler.GetICP))                         // 查询备案
	api.GET("/tool/ping", ginHelper.Handle(handler.Ping))                          // TODO ping

}

func Data() {
	data := Router.Group("/data")
	data.GET("/allurl/:host", ginHelper.Handle(handler.WebsiteAllUrl))
	data.GET("/all/website", ginHelper.Handle(handler.AllWebsite))
}

func Test() {
	// test := Router.Group("/test")
	// test.GET("/css", ginHelper.Handle(handler.GetCss)) // 获取css js
}

// AuthPG 权限验证中间件
func AuthPG() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(constname.UserToken)
		j := utils.NewJWT(conf.Conf.Default.Jwt.Secret, conf.Conf.Default.Jwt.Expire)
		if err := j.ParseToken(token); err == nil {
			c.Next()
			return
		}
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}
}

// AuthAPI 权限验证中间件
func AuthAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(constname.UserToken)
		j := utils.NewJWT(conf.Conf.Default.Jwt.Secret, conf.Conf.Default.Jwt.Expire)
		if err := j.ParseToken(token); err == nil {
			c.Next()
			return
		}
		ginHelper.AuthErrorOut(c)
		c.Abort()
		return

	}
}
