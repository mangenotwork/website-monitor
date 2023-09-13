package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/ginHelper"
	"github.com/mangenotwork/common/utils"
	"net/http"
	"website-monitor/master/constname"
)

func ginH(h gin.H) gin.H {
	h["Title"] = conf.Conf.Default.App.Name
	h["TimeStamp"] = constname.TimeStamp
	return h
}

func NotFond(c *gin.Context) {
	// 实现内部重定向
	c.HTML(
		http.StatusOK,
		"notfond.html",
		ginH(gin.H{}),
	)
}

func ErrPage(c *gin.Context, err error) {
	c.HTML(
		http.StatusOK,
		"err.html",
		ginH(gin.H{
			"err": err.Error(),
		}),
	)
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
}

func HomePage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"home.html",
		ginH(gin.H{
			"nav": "home",
		}),
	)
}
