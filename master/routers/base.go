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
	api.GET("/tool/certificate", ginHelper.Handle(handler.GetSSLCertificate)) // 获取证书
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
