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

	// mail
	api.GET("/mail/init", ginHelper.Handle(handler.MailInit))          // 是否设置邮件
	api.POST("/mail/conf", ginHelper.Handle(handler.MailConf))         // 设置邮件配置
	api.GET("/mail/info", ginHelper.Handle(handler.MailInfo))          // 获取邮件配置信息
	api.POST("/mail/sendTest", ginHelper.Handle(handler.MailSendTest)) // 测试发生邮件

	// tool
	api.GET("/tool/certificate", ginHelper.Handle(handler.GetSSLCertificate))      // 获取证书
	api.GET("/tool/nsLookUp", ginHelper.Handle(handler.DNSLookUp))                 // 查询dns
	api.GET("/tool/nsLookUp/all", ginHelper.Handle(handler.DNSLookUpAll))          // 查询dns
	api.GET("/tool/whois", ginHelper.Handle(handler.Whois))                        // Whois查询
	api.GET("/tool/ip", ginHelper.Handle(handler.IPInfo))                          // ip信息查询
	api.GET("/tool/myIP", ginHelper.Handle(handler.MyIPInfo))                      // 本机ip信息
	api.GET("/tool/website/tdki", ginHelper.Handle(handler.GetWebSiteTDKI))        // 获取网站的T, D, K, 图标
	api.GET("/tool/website/collectInfo", ginHelper.Handle(handler.CollectWebSite)) // 采集网站信息
	api.GET("/tool/icp", ginHelper.Handle(handler.GetICP))                         // 查询备案

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
