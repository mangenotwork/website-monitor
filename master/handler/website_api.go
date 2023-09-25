package handler

import (
	"time"

	"website-monitor/master/constname"
	"website-monitor/master/dao"
	"website-monitor/master/entity"

	"github.com/mangenotwork/common/ginHelper"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
)

type WebsiteAddParam struct {
	Host                     string `json:"host"`
	Notes                    string `json:"notes"`
	MonitorRate              int64  `json:"monitorRate"`
	ContrastUrl              string `json:"contrastUrl"`
	ContrastTime             int64  `json:"contrastTime"`
	Ping                     string `json:"ping"`
	PingTime                 int64  `json:"pingTime"`
	UriDepth                 int64  `json:"uriDepth"`
	ScanRate                 int64  `json:"scanRate"`
	ScanBadLink              bool   `json:"scanBadLink"`
	ScanTDK                  bool   `json:"scanTDK"`
	ScanExtLinks             bool   `json:"scanExtLinks"`
	WebsiteSlowResponseTime  int64  `json:"websiteSlowResponseTime"`
	WebsiteSlowResponseCount int64  `json:"websiteSlowResponseCount"`
	SSLCertificateExpire     int64  `json:"SSLCertificateExpire"`
}

func WebsiteAdd(c *ginHelper.GinCtx) {
	param := &WebsiteAddParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		c.APIOutPutError(err, "参数或参数类型错误")
		return
	}
	if len(param.Host) < 1 {
		c.APIOutPutError(nil, "参数错误: host不能为空")
		return
	}
	if param.MonitorRate < 1 {
		param.MonitorRate = constname.DefaultMonitorRate
	}
	if len(param.ContrastUrl) < 1 {
		param.ContrastUrl = constname.DefaultContrastUrl
	}
	if param.ContrastTime < 1 {
		param.ContrastTime = constname.DefaultContrastTime
	}
	if len(param.Ping) < 1 {
		param.Ping = constname.DefaultPing
	}
	if param.PingTime < 1 {
		param.PingTime = constname.DefaultPingTime
	}
	if param.UriDepth < 1 {
		param.UriDepth = constname.DefaultUriDepth
	}
	if param.ScanRate < 1 {
		param.ScanRate = constname.DefaultScanRate
	}
	if param.WebsiteSlowResponseTime < 100 {
		c.APIOutPutError(nil, "网站响应慢不能小于100ms")
		return
	}
	if param.WebsiteSlowResponseCount < 1 {
		param.WebsiteSlowResponseCount = constname.DefaultWebsiteSlowResponseCount
	}
	if param.SSLCertificateExpire < 1 {
		param.SSLCertificateExpire = constname.DefaultSSLCertificateExpire
	}

	log.Info("param = ", param)
	website := &entity.Website{
		Host:         param.Host,
		MonitorRate:  param.MonitorRate,
		ContrastUrl:  param.ContrastUrl,
		ContrastTime: param.ContrastTime,
		Ping:         param.Ping,
		PingTime:     param.PingTime,
		Notes:        param.Notes,
		Created:      time.Now().Unix(),
	}
	alarmRule := &entity.WebsiteAlarmRule{
		Host:                     param.Host,
		WebsiteSlowResponseTime:  param.WebsiteSlowResponseTime,
		WebsiteSlowResponseCount: param.WebsiteSlowResponseCount,
		SSLCertificateExpire:     param.SSLCertificateExpire,
		NotTDK:                   param.ScanTDK,
		BadLink:                  param.ScanBadLink,
		ExtLinkChange:            param.ScanExtLinks,
	}
	scan := &entity.WebsiteScanCheckUp{
		Host:         param.Host,
		ScanDepth:    param.UriDepth,
		ScanRate:     param.ScanRate,
		ScanTDK:      param.ScanTDK,
		ScanBadLink:  param.ScanBadLink,
		ScanExtLinks: param.ScanExtLinks,
	}
	err = dao.NewWebsite().Add(website, alarmRule, scan)
	if err != nil {
		c.APIOutPutError(err, "创建失败, err = "+err.Error())
		return
	}
	c.APIOutPut("创建成功", "创建成功")
	return
}

func WebsiteList(c *ginHelper.GinCtx) {
	data, _, err := dao.NewWebsite().SelectList()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(data, "")
	return
}

func WebsiteDelete(c *ginHelper.GinCtx) {
	host := c.Param("host")
	err := dao.NewWebsite().Del(host)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut("ok", "成功删除")
	return
}

type WebsiteInfoOutPut struct {
	Base        WebsiteOutPut              `json:"base"`
	Info        *entity.WebsiteInfo        `json:"info"`
	AlarmRule   *entity.WebsiteAlarmRule   `json:"alarmRule"`
	ScanCheckUp *entity.WebsiteScanCheckUp `json:"scanCheckUp"`
}

type WebsiteOutPut struct {
	*entity.Website
	Date string `json:"date"`
}

func WebsiteInfo(c *ginHelper.GinCtx) {
	host := c.Param("host")
	website := dao.NewWebsite()
	output := &WebsiteInfoOutPut{}
	var err error
	base, err := website.Select(host)
	output.Base = WebsiteOutPut{base, utils.Timestamp2Date(base.Created)}
	output.Info, err = website.GetInfo(host)
	output.AlarmRule, err = website.GetAlarmRule(host)
	output.ScanCheckUp, err = website.GetScanCheckUp(host)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(output, "")
	return
}

func WebsiteInfoRefresh(c *ginHelper.GinCtx) {
	host := c.GetQuery("host")
	hostId := c.GetQuery("id")
	log.Info("WebsiteInfoRefresh")
	err := dao.NewWebsite().SaveCollectInfo(host, hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut("ok", "")
	return
}

func WebsiteUrls(c *ginHelper.GinCtx) {
	hostId := c.Param("host")
	data, err := dao.NewWebsite().GetWebSiteUrl(hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(data, "")
	return
}

func WebsiteEdit(c *ginHelper.GinCtx) {

}

func WebsiteChart(c *ginHelper.GinCtx) {

}

func WebsiteAlertList(c *ginHelper.GinCtx) {

}

func WebsiteAlertDel(c *ginHelper.GinCtx) {

}
