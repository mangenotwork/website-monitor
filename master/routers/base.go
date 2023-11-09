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
	pg.GET("/alert", handler.AlertPage)
	pg.GET("/tool", handler.ToolPage)
	pg.GET("/test/api", handler.TestAPIPage)
	pg.GET("/test/stress", handler.TestStressPage)
	pg.GET("/test/penetration", handler.TestPenetrationPage)
	pg.GET("/operation", handler.OperationPage)
	pg.GET("/instructions", handler.InstructionsPage)
	pg.GET("/requester", handler.RequesterPage)
}

func API() {
	api := Router.Group("/api").Use(AuthAPI())
	api.GET("/out", handler.Out)

	// website
	api.POST("/website/add", ginHelper.Handle(handler.WebsiteAdd))                       // 创建网站监测
	api.GET("/website/list", ginHelper.Handle(handler.WebsiteList))                      // 监测网站列表
	api.GET("/website/conf/:hostId", ginHelper.Handle(handler.WebsiteConf))              // 监测网站的监测配置信息
	api.GET("/website/info/:hostId", ginHelper.Handle(handler.WebsiteInfo))              // 监测网站详情
	api.GET("/website/info/refresh", ginHelper.Handle(handler.WebsiteInfoRefresh))       // 监测网站详情刷新
	api.GET("/website/delete/:hostId", ginHelper.Handle(handler.WebsiteDelete))          // 删除网站监测
	api.GET("/website/urls/:hostId", ginHelper.Handle(handler.WebsiteUrls))              // 监测网站采集到url
	api.POST("/website/edit/:hostId", ginHelper.Handle(handler.WebsiteEdit))             // 监测设置
	api.GET("/website/chart/:hostId", ginHelper.Handle(handler.WebsiteChart))            // 图表
	api.GET("/website/log/:hostId", ginHelper.Handle(handler.MonitorLog))                // 获取监测日志
	api.GET("/website/log/list/:hostId", ginHelper.Handle(handler.MonitorLogList))       // 获取监测日志列表
	api.GET("/website/log/upload/:hostId", ginHelper.Handle(handler.MonitorLogUpload))   // 日志文件下载
	api.POST("/website/point/add/:hostId", ginHelper.Handle(handler.WebsitePointAdd))    // 添加监测点
	api.GET("/website/point/list/:hostId", ginHelper.Handle(handler.WebsitePointList))   // 获取监测点
	api.POST("/website/point/del/:hostId", ginHelper.Handle(handler.WebsitePointDel))    // 删除监测点
	api.GET("/website/point/clear/:hostId", ginHelper.Handle(handler.WebsitePointClear)) // 清空监测点

	// mail
	api.GET("/mail/init", ginHelper.Handle(handler.MailInit))          // 是否设置邮件
	api.POST("/mail/conf", ginHelper.Handle(handler.MailConf))         // 设置邮件配置
	api.GET("/mail/info", ginHelper.Handle(handler.MailInfo))          // 获取邮件配置信息
	api.POST("/mail/sendTest", ginHelper.Handle(handler.MailSendTest)) // 测试发生邮件

	// alert
	api.GET("/alert/list", ginHelper.Handle(handler.AlertList))               // 报警列表
	api.GET("/alert/wbesite/:hostId", ginHelper.Handle(handler.AlertWebsite)) // 指定网站的报警信息
	api.GET("/alert/read/:id", ginHelper.Handle(handler.AlertRead))           // 报警信息已读
	api.GET("/alert/info/:id", ginHelper.Handle(handler.AlertInfo))           // 报警消息详细信息
	api.GET("/alert/del/:id", ginHelper.Handle(handler.AlertDel))             // 报警消息删除
	api.GET("/alert/clear/:hostId", ginHelper.Handle(handler.AlertClear))     // 指定网站的报警信息清空
	api.GET("/alert/all/clear", ginHelper.Handle(handler.AlertAllClear))      // 指定网站的报警信息清空

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
	api.GET("/tool/ping", ginHelper.Handle(handler.Ping))                          // ping
	// TODO 网站宿主ip端口扫描

	// monitor
	api.GET("/monitor/list", ginHelper.Handle(handler.MonitorList))     // 监测器列表，在线情况
	api.GET("/monitor/ipaddr", ginHelper.Handle(handler.MonitorIPAddr)) // 测试获取ip地址属地信息

	// requester
	api.GET("/requester/create/tab", ginHelper.Handle(handler.RequesterCreateTab))                // 创建一个新请求标签页
	api.GET("/requester/close/tab/:reqId", ginHelper.Handle(handler.RequesterCloseTab))           // 删除一个标签页
	api.POST("/requester/execute", ginHelper.Handle(handler.RequesterExecute))                    // TODO 请求调试执行
	api.GET("/requester/data/:reqId", ginHelper.Handle(handler.RequesterGetData))                 // 获取指定请求数据
	api.GET("/requester/list", ginHelper.Handle(handler.RequesterList))                           // 请求调试列表
	api.GET("/requester/history/list", ginHelper.Handle(handler.RequesterHistoryList))            // 请求调试历史列表
	api.GET("/requester/history/delete/:reqId", ginHelper.Handle(handler.RequesterHistoryDelete)) // 请求调试历史删除
	api.GET("/requester/history/creat", ginHelper.Handle(handler.RequesterHistoryDelete))         // TODO 请求调试历史删除
	api.GET("/requester/dir/creat", ginHelper.Handle(handler.RequesterDirCreat))                  // TODO 请求目录-创建
	api.GET("/requester/dir/list", ginHelper.Handle(handler.RequesterDirList))                    // TODO 请求目录-列表
	api.GET("/requester/dir/join", ginHelper.Handle(handler.RequesterDirJoin))                    // TODO 请求目录-加入
	api.POST("/requester/global/header", ginHelper.Handle(handler.RequesterGlobalHeaderSet))      // 设置全局Header
	api.GET("/requester/global/header", ginHelper.Handle(handler.RequesterGlobalHeaderGet))       // 获取全局Header
	api.GET("/requester/global/header/del", ginHelper.Handle(handler.RequesterGlobalHeaderDel))   // 删除指定Header

}

func Data() {
	// 提供监测器获取数据
	data := Router.Group("/data")
	data.GET("/allurl/:hostId", ginHelper.Handle(handler.WebsiteAllUrl))           // 获取网站下的所有url
	data.GET("/all/website", ginHelper.Handle(handler.AllWebsite))                 // 获取所有需要监测的网站
	data.GET("/website/point/:hostId", ginHelper.Handle(handler.WebsitePointList)) // 获取网站监测点
	data.GET("/website/:hostId", ginHelper.Handle(handler.GetWebsiteData))         // 获取指定网站信息
}

func Test() {
	test := Router.Group("/test")
	test.GET("/NoticeUpdateWebsite", ginHelper.Handle(handler.NoticeUpdateWebsiteTest))
	test.GET("/NoticeUpdateWebsiteAllUrl", ginHelper.Handle(handler.NoticeUpdateWebsiteAllUrlTest))
	test.GET("/NoticeUpdateWebsitePoint", ginHelper.Handle(handler.NoticeUpdateWebsitePointTest))
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
