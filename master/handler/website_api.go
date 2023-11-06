package handler

import (
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"sort"
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

type WebsiteListOutItem struct {
	*entity.Website
	AlertCount int `json:"alertCount"`
	State      int `json:"state"` // 0:未执行(没有监测器连接)  1:监测中  2:有报警
}

func WebsiteList(c *ginHelper.GinCtx) {
	data := make([]*WebsiteListOutItem, 0)
	websiteList, _, err := dao.NewWebsite().SelectList()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	// 监测器状态
	isMonitor := false
	monitor := dao.GetClientList()
	for _, v := range monitor {
		if v.Online {
			isMonitor = true
			break
		}
	}
	// 报警数量获取
	for _, v := range websiteList {
		alertList, _ := dao.NewAlert().GetWebsiteAlertIDList(v.HostID)
		alertLen := len(alertList)
		state := 0
		if isMonitor {
			state = 1
		}
		if alertLen > 0 {
			state = 2
		}
		data = append(data, &WebsiteListOutItem{
			v,
			alertLen,
			state,
		})
	}
	sort.Slice(data, func(i, j int) bool {
		if data[i].Created > data[j].Created {
			return true
		}
		return false
	})
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
	obj := dao.NewWebsite()
	data, err := obj.GetWebSiteUrl(hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	website, _ := obj.Select(hostId)
	if len(data.AllUri) == 1 && data.AllUri[0] == website.Host {
		c.APIOutPut([]string{}, "")
		return
	}
	c.APIOutPut(data.AllUri, "")
	return
}

func AllWebsite(c *ginHelper.GinCtx) {
	data := make([]*WebsiteDataOut, 0)
	obj := dao.NewWebsite()
	websiteList, _, err := obj.SelectList()
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	for _, v := range websiteList {
		rule, rErr := obj.GetAlarmRule(v.HostID)
		if rErr != nil {
			continue
		}
		data = append(data, &WebsiteDataOut{v, rule.WebsiteSlowResponseTime})
	}
	c.APIOutPut(data, "")
	return
}

type WebsiteDataOut struct {
	*entity.Website
	WebsiteSlowResponseTime int64
}

func GetWebsiteData(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	obj := dao.NewWebsite()
	website, err := obj.Select(hostId)
	rule, err := obj.GetAlarmRule(hostId)
	if err != nil {
		c.APIOutPutError(err, err.Error())
		return
	}
	data := &WebsiteDataOut{
		website,
		rule.WebsiteSlowResponseTime,
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

type AlertListOut struct {
	List  []*entity.AlertData `json:"list"`
	Count int                 `json:"count"`
}

func AlertList(c *ginHelper.GinCtx) {
	list, err := dao.NewAlert().GetList()
	if err != nil {
		log.Error(err)
		c.APIOutPutError(nil, err.Error())
		return
	}
	data := &AlertListOut{
		List:  list,
		Count: len(list),
	}
	c.APIOutPut(data, "")
	return
}

func AlertWebsite(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	list, err := dao.NewAlert().GetAtWebsite(hostId)
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	data := &AlertListOut{
		List:  list,
		Count: len(list),
	}
	c.APIOutPut(data, "")
	return
}

func AlertRead(c *ginHelper.GinCtx) {
	id := c.Param("id")
	err := dao.NewAlert().Read(id)
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	c.APIOutPut("标记成功", "标记成功")
	return
}

func AlertInfo(c *ginHelper.GinCtx) {
	id := c.Param("id")
	data, err := dao.NewAlert().Get(id)
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	c.APIOutPut(data, "")
	return
}

func AlertDel(c *ginHelper.GinCtx) {
	id := c.Param("id")
	err := dao.NewAlert().Del(id)
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	c.APIOutPut("成功", "成功")
	return
}

func AlertClear(c *ginHelper.GinCtx) {
	hostId := c.Param("hostId")
	err := dao.NewAlert().Clear(hostId)
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	c.APIOutPut("成功", "成功")
	return
}

func AlertAllClear(c *ginHelper.GinCtx) {
	err := dao.NewAlert().ClearAll()
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	c.APIOutPut("成功", "成功")
	return
}

type RequesterExecuteParam struct {
	Name         string         `json:"name"`         // api name
	Note         string         `json:"note"`         // api note
	Method       string         `json:"method"`       // 请求类型
	Url          string         `json:"url"`          // 请求url
	Header       map[string]any `json:"header"`       // 请求header
	BodyType     string         `json:"bodyType"`     // 请求body type
	BodyJson     string         `json:"bodyJson"`     // body json
	BodyFromData map[string]any `json:"bodyFromData"` // body from-data
	BodyXWWWFrom map[string]any `json:"bodyXWWWFrom"` // body x-www-from
	BodyXml      string         `json:"bodyXml"`      // body xml
	BodyText     string         `json:"bodyText"`     // body text
}

func RequesterExecute(c *ginHelper.GinCtx) {
	param := &RequesterExecuteParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	if !isMethod(param.Method) {
		c.APIOutPutError(nil, "未知的请求类型")
		return
	}
	if len(param.Url) < 1 {
		c.APIOutPutError(nil, "请求地址为空")
		return
	}
	if len(param.Name) == 0 {
		param.Name = "新建请求"
	}
	log.Info("param = ", param)
}

func isMethod(method string) bool {
	rse := false
	for _, v := range []string{"get", "post", "put", "delete", "options", "head"} {
		if method == v {
			rse = true
			break
		}
	}
	return rse
}

func RequesterList(c *ginHelper.GinCtx) {

}

func RequesterHistoryList(c *ginHelper.GinCtx) {

}

func RequesterHistoryDelete(c *ginHelper.GinCtx) {

}

func RequesterDirCreat(c *ginHelper.GinCtx) {

}

func RequesterDirList(c *ginHelper.GinCtx) {

}

func RequesterDirJoin(c *ginHelper.GinCtx) {

}

type RequesterGlobalHeaderSetParam struct {
	List []*entity.RequesterGlobalHeader `json:"list"`
}

func RequesterGlobalHeaderSet(c *ginHelper.GinCtx) {
	param := &RequesterGlobalHeaderSetParam{}
	err := c.GetPostArgs(param)
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	log.Info("param = ", param)
	err = dao.NewRequestTool().SetGlobalHeader(param.List)
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	c.APIOutPut("成功", "成功")
	return
}

func RequesterGlobalHeaderGet(c *ginHelper.GinCtx) {
	data, err := dao.NewRequestTool().GetGlobalHeader()
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	c.APIOutPut(data, "成功")
	return
}

func RequesterGlobalHeaderDel(c *ginHelper.GinCtx) {
	key := c.Query("key")
	err := dao.NewRequestTool().DelGlobalHeader(key)
	if err != nil {
		c.APIOutPutError(nil, err.Error())
		return
	}
	c.APIOutPut("成功", "成功")
	return
}
