package handler

import (
	"net/http"
	"time"

	"website-monitor/master/constname"
	"website-monitor/master/dao"
	"website-monitor/master/entity"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/ginHelper"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
)

func Login(c *gin.Context) {
	user := c.PostForm("user")
	password := c.PostForm("password")
	for _, v := range conf.Conf.Default.User {
		if user == v.Name && password == v.Password {
			j := utils.NewJWT(conf.Conf.Default.Jwt.Secret, conf.Conf.Default.Jwt.Expire)
			j.AddClaims("name", user)
			token, tokenErr := j.Token()
			if tokenErr != nil {
				log.Error("生产token错误， err = ", tokenErr)
			}
			c.SetCookie(constname.UserToken, token, constname.TokenExpires,
				"/", "", false, true)
			c.Redirect(http.StatusFound, "/home")
			return
		}
	}
	c.HTML(200, "err.html", gin.H{
		"Title": conf.Conf.Default.App.Name,
		"err":   "账号或密码错误",
	})
	return
}

func Out(c *gin.Context) {
	c.SetCookie("sign", "", 60*60*24*7, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func MailInit(c *ginHelper.GinCtx) {
	data := dao.NewMail().IsMail()
	c.APIOutPut(data, "")
	return
}

func mailSet(c *ginHelper.GinCtx) error {
	param := &entity.Mail{}
	err := c.GetPostArgs(param)
	if err != nil {
		log.Error(err)
		return err
	}
	mailDao := dao.NewMail()
	err = mailDao.Check(param)
	if err != nil {
		return err
	}
	err = mailDao.SetMail(param)
	if err != nil {
		return err
	}
	return nil
}

func MailConf(c *ginHelper.GinCtx) {
	err := mailSet(c)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut("设置成功", "设置成功!")
	return
}

func MailInfo(c *ginHelper.GinCtx) {
	data, err := dao.NewMail().GetMail()
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	c.APIOutPut(data, "")
	return
}

func MailSendTest(c *ginHelper.GinCtx) {
	err := mailSet(c)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	title := "Website Monitor 邮件通知测验"
	body := `<p>您好欢迎使用Website Monitor，这是一封邮件通知测验的邮件，你收到此邮件说明监测平台通知配置成功!</p>` +
		`<p> 开源地址: </p>` +
		`<p><a herf="https://github.com/mangenotwork/website-monitor">https://github.com/mangenotwork/website-monitor</a></p>` +
		`<p>ManGe : ` + time.Now().String() + `</p>`
	dao.NewMail().Send(title, body)
	c.APIOutPut("", "测试邮件已发送请注意查收!")
	return
}

func GetSSLCertificate(c *ginHelper.GinCtx) {
	caseUrl := c.GetQuery("url")
	log.Info(caseUrl)
	data, _ := dao.GetCertificateInfo(caseUrl)
	c.APIOutPut(data, "")
	return
}

type DNSLookUpAllOut struct {
	List []*entity.DNSInfo `json:"list"`
	IPs  []string          `json:"ips"`
}

func DNSLookUp(c *ginHelper.GinCtx) {
	host := c.GetQuery("host")
	data := dao.NsLookUpLocal(host)
	c.APIOutPut(data, "")
	return
}

func DNSLookUpAll(c *ginHelper.GinCtx) {
	host := c.GetQuery("host")
	list, allIP := dao.NsLookUpAll(host)
	c.APIOutPut(&DNSLookUpAllOut{
		List: list,
		IPs:  allIP,
	}, "")
	return
}

func Whois(c *ginHelper.GinCtx) {
	host := c.GetQuery("host")
	data := dao.Whois(host)
	c.APIOutPut(data, "")
	return
}

func IPInfo(c *ginHelper.GinCtx) {
	ip := c.GetQuery("ip")
	data := dao.GetIP(ip)
	c.APIOutPut(data, "")
	return
}

func MyIPInfo(c *ginHelper.GinCtx) {
	data := dao.GetMyIP()
	c.APIOutPut(data, "")
	return
}

func GetWebSiteTDKI(c *ginHelper.GinCtx) {
	url := c.GetQuery("url")
	data := dao.NewWebsite().CollectTDK(url)
	c.APIOutPut(data, "")
	return
}
