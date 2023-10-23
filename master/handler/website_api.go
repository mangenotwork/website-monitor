package handler

import (
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"time"
	"website-monitor/master/business"

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

func analysisWebsiteAddParam(c *ginHelper.GinCtx) (*WebsiteAddParam, error) {
	param := &WebsiteAddParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		return param, err
	}
	if len(param.Host) < 1 {
		return param, fmt.Errorf("参数错误: host不能为空")
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
		return param, fmt.Errorf("网站响应慢不能小于100ms")
	}
	if param.WebsiteSlowResponseCount < 1 {
		param.WebsiteSlowResponseCount = constname.DefaultWebsiteSlowResponseCount
	}
	if param.SSLCertificateExpire < 1 {
		param.SSLCertificateExpire = constname.DefaultSSLCertificateExpire
	}
	return param, nil
}

func WebsiteAdd(c *ginHelper.GinCtx) {
	param, err := analysisWebsiteAddParam(c)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
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
	// TODO 获取监测结果数据
	data, _, err := dao.NewWebsite().SelectList()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(data, "")
	return
}

type WebsiteConfOut struct {
	Base        *entity.Website            `json:"base"`
	AlarmRule   *entity.WebsiteAlarmRule   `json:"alarmRule"`
	ScanCheckUp *entity.WebsiteScanCheckUp `json:"scanCheckUp"`
}

func WebsiteConf(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	website := dao.NewWebsite()
	base, err := website.Select(hostId)
	log.Info("base = ", base)
	alertRule, err := website.GetConfAlarmRule(hostId)
	checkUp, err := website.GetConfScanCheckUp(hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	out := &WebsiteConfOut{base, alertRule, checkUp}
	c.APIOutPut(out, "")
	return
}

func WebsiteDelete(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	err := dao.NewWebsite().Del(hostId)
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
	hostId := c.Param("hostId")
	website := dao.NewWebsite()
	output := &WebsiteInfoOutPut{}
	var err error
	base, err := website.Select(hostId)
	output.Base = WebsiteOutPut{base, utils.Timestamp2Date(base.Created)}
	output.Info, err = website.GetInfo(hostId)
	output.AlarmRule, err = website.GetAlarmRule(hostId)
	output.ScanCheckUp, err = website.GetScanCheckUp(hostId)
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
	hostId := c.Param("hostId")
	data, err := dao.NewWebsite().GetWebSiteUrl(hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(data, "")
	return
}

func WebsiteAllUrl(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	data, err := dao.NewWebsite().GetWebSiteUrl(hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(data.AllUri, "")
	return
}

func AllWebsite(c *ginHelper.GinCtx) {
	data, _, err := dao.NewWebsite().SelectList()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(data, "")
	return
}

func GetWebsiteData(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	data, err := dao.NewWebsite().Select(hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(data, "")
	return
}

func WebsiteEdit(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	param, err := analysisWebsiteAddParam(c)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	log.Info("param = ", param)
	website := &entity.Website{
		Host:         param.Host,
		HostID:       hostId,
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
		HostID:                   hostId,
		WebsiteSlowResponseTime:  param.WebsiteSlowResponseTime,
		WebsiteSlowResponseCount: param.WebsiteSlowResponseCount,
		SSLCertificateExpire:     param.SSLCertificateExpire,
		NotTDK:                   param.ScanTDK,
		BadLink:                  param.ScanBadLink,
		ExtLinkChange:            param.ScanExtLinks,
	}
	scan := &entity.WebsiteScanCheckUp{
		Host:         param.Host,
		HostID:       hostId,
		ScanDepth:    param.UriDepth,
		ScanRate:     param.ScanRate,
		ScanTDK:      param.ScanTDK,
		ScanBadLink:  param.ScanBadLink,
		ScanExtLinks: param.ScanExtLinks,
	}
	err = dao.NewWebsite().Edit(website, alarmRule, scan)
	if err != nil {
		c.APIOutPutError(err, "修改监测配置失败, err = "+err.Error())
		return
	}
	c.APIOutPut("修改监测配置成功", "修改监测配置成功")
	return
}

func WebsiteChart(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	day := c.GetQuery("day")
	if len(day) < 1 {
		day = utils.NowDateLayout(constname.DayLayout)
	}
	uri := c.GetQuery("uri") // Health:根URI,健康URI  Random:随机URI  Point:监测点URI
	data, err := dao.NewMonitorLogDao().ReadAll(hostId, day)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	out := make([][]int64, 0)
	for _, v := range data {
		if len(uri) > 0 { // 指定类型
			if v.UriType == uri {
				item := make([]int64, 0)
				item = append(item, utils.Date2Timestamp(v.Time)*1000)
				item = append(item, v.UriMs)
				out = append(out, item)
			}
		} else {
			item := make([]int64, 0)
			item = append(item, utils.Date2Timestamp(v.Time)*1000)
			item = append(item, v.UriMs)
			out = append(out, item)
		}
	}
	c.APIOutPut(out, "")
	return
}

func WebsiteAlertList(c *ginHelper.GinCtx) {

}

func WebsiteAlertDel(c *ginHelper.GinCtx) {

}

func MonitorLog(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	day := c.GetQuery("day")
	data := dao.NewMonitorLogDao().ReadLog(hostId, day)
	c.APIOutPut(data, "")
	return
}

func MonitorLogList(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	data, err := dao.NewMonitorLogDao().LogListDay(hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(data, "")
	return
}

func MonitorLogUpload(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	day := c.GetQuery("day")
	logPath, err := dao.NewMonitorLogDao().Upload(hostId, day)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment; filename=%s", fmt.Sprintf("%s.log", day)))
	c.Writer.Header().Add("Content-Type", "application/text/plain")
	c.File(logPath)
	return
}

type WebsitePointParam struct {
	Uri string `json:"uri"`
}

func WebsitePointAdd(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	log.Info("hostId = ", hostId)
	param := &WebsitePointParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		c.APIOutPutError(err, "参数或参数类型错误")
		return
	}
	ctx, _ := gt.Get(param.Uri)
	if business.AlertRuleCode(ctx.StateCode) {
		c.APIOutPutError(nil, fmt.Sprintf("%s请求失败，状态码:%d", param.Uri, ctx.StateCode))
		return
	}
	err = dao.NewWebsite().SetPoint(hostId, param.Uri)
	if err != nil {
		c.APIOutPutError(err, "添加监测点失败:"+err.Error())
		return
	}
	c.APIOutPut("", "添加成功")
	return
}

func WebsitePointList(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	data, err := dao.NewWebsite().GetPoint(hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut(data, "成功")
	return
}

func WebsitePointDel(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	log.Info("hostId = ", hostId)
	param := &WebsitePointParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		c.APIOutPutError(err, "参数或参数类型错误")
		return
	}
	err = dao.NewWebsite().DelPoint(hostId, param.Uri)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut("", "成功")
	return
}

func WebsitePointClear(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	err := dao.NewWebsite().ClearPoint(hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	c.APIOutPut("", "成功")
	return
}

func MonitorList(c *ginHelper.GinCtx) {
	data := dao.GetClientList()
	c.APIOutPut(data, "成功")
	return
}

func MonitorIPAddr(c *ginHelper.GinCtx) {
	dao.GetClientList2()
	c.APIOutPut("ok", "成功")
	return
}
