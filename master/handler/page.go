package handler

import (
	"net/http"

	"website-monitor/master/constname"

	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/ginHelper"
	"github.com/mangenotwork/common/utils"
)

func ginH(h gin.H) gin.H {
	h["Title"] = conf.Conf.Default.App.Name
	h["TimeStamp"] = constname.TimeStamp
	return h
}

func NotFond(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"not_fond.html",
		ginH(gin.H{}),
	)
	return
}

func ErrPage(c *gin.Context) {
	msg := c.Query("msg")
	c.HTML(
		http.StatusOK,
		"err.html",
		ginH(gin.H{
			"err":       msg,
			"returnUrl": "/",
		}),
	)
	return
}

func LoginPage(c *gin.Context) {
	token, _ := c.Cookie(constname.UserToken)

	if token != "" {

		j := utils.NewJWT(conf.Conf.Default.Jwt.Secret, conf.Conf.Default.Jwt.Expire)
		if err := j.ParseToken(token); err == nil {
			c.Redirect(http.StatusFound, "/home")
			return
		}

	}

	c.HTML(
		http.StatusOK,
		"login.html",
		ginH(gin.H{
			"csrf": ginHelper.FormSetCSRF(c.Request),
		}),
	)
	return
}

func HomePage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"home.html",
		ginH(gin.H{
			"nav": "home",
		}),
	)
	return
}

func MonitorPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"monitor.html",
		ginH(gin.H{
			"nav": "monitor",
		}),
	)
	return
}

func AlertPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"alert.html",
		ginH(gin.H{
			"nav": "alert",
		}),
	)
	return
}

func ToolPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"tool.html",
		ginH(gin.H{
			"nav": "tool",
		}),
	)
	return
}

func TestAPIPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"test_api.html",
		ginH(gin.H{
			"nav": "api-test",
		}),
	)
	return
}

func TestStressPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"test_stress.html",
		ginH(gin.H{
			"nav": "stress-test",
		}),
	)
	return
}

func TestPenetrationPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"test_penetration.html",
		ginH(gin.H{
			"nav": "penetration-test",
		}),
	)
	return
}

func OperationPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"operation.html",
		ginH(gin.H{
			"nav": "operation",
		}),
	)
	return
}

func InstructionsPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"instructions.html",
		ginH(gin.H{
			"nav": "instructions",
		}),
	)
	return
}

func RequesterPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"requester.html",
		ginH(gin.H{
			"nav": "requester",
		}),
	)
	return
}
