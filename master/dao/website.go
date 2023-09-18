package dao

import (
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"website-monitor/master/entity"
)

type WebsiteEr interface {
	Add(data *entity.Website, alarmRule *entity.WebsiteAlarmRule, scan *entity.WebsiteScanCheckUp) error
	Del()
	Update()
	SelectList()
	Select()

	// 报警列表
	AlertList()

	// 获取采集的网站Url
	GetWebSiteUrl()

	// 设置指定监测Url点
	SetUrlPoint()

	// 获取监测Url点
	GetUrlPoint()

	// 采集网站页面基础信息
	CollectTDK(host string) *entity.TDKI

	// 采集网站DNS信息
	CollectDNS()

	// 采集网站IP和属地
	CollectIPAddr()

	// 采集网站证书信息
	CollectSSLCertificateInfo()

	// 采集网站Whois信息
	CollectWhois()
}

func NewWebsite() WebsiteEr {
	return new(websiteDao)
}

type websiteDao struct{}

func (w *websiteDao) Add(data *entity.Website, rule *entity.WebsiteAlarmRule, scan *entity.WebsiteScanCheckUp) error {
	// 检查站点是否可访问
	ctx, err := gt.Get(data.Host)
	if err != nil {
		return err
	}
	if !InspectCode(ctx.StateCode) {
		return fmt.Errorf("网站:%s 请求状态码为 %d , 无法添加。", data.Host, ctx.StateCode)
	}
	err = w.addWebsite(data)
	if err != nil {
		return err
	}
	// 保存报警规则信息
	rule.Host = data.Host
	err = w.addWebsiteAlarmRule(rule)
	if err != nil {
		return err
	}
	// 保存扫描规则信息
	scan.Host = data.Host
	err = w.addWebsiteScanCheckUp(scan)
	if err != nil {
		return err
	}
	// TODO 异步执行扫描

	// TODO... 异步获取网站信息

	return err
}

func (w *websiteDao) addWebsite(data *entity.Website) error {
	return DB.Set(WebSiteTable, data.Host, data)
}

func (w *websiteDao) addWebsiteAlarmRule(alarmRule *entity.WebsiteAlarmRule) error {
	if alarmRule.Host == "" {
		return fmt.Errorf("没有host")
	}
	return DB.Set(WebsiteAlarmRuleTable, alarmRule.Host, alarmRule)
}

func (w *websiteDao) addWebsiteScanCheckUp(scan *entity.WebsiteScanCheckUp) error {
	if scan.Host == "" {
		return fmt.Errorf("没有host")
	}
	return DB.Set(WebsiteScanCheckUpTable, scan.Host, scan)
}

func (w *websiteDao) Del() {

}

func (w *websiteDao) Update() {

}

func (w *websiteDao) SelectList() {

}

func (w *websiteDao) Select() {

}

func (w *websiteDao) CollectTDK(host string) *entity.TDKI {
	tdki := &entity.TDKI{
		Icon: host + "/favicon.ico",
	}
	ctx, _ := gt.Get(host)
	title := gt.RegHtmlTitleTxt(ctx.Html)
	if len(title) > 0 {
		tdki.Title = title[0]
	}
	keyword := gt.RegHtmlKeywordTxt(ctx.Html)
	if len(keyword) > 0 {
		tdki.Keywords = keyword[0]
	}
	description := gt.RegHtmlDescriptionTxt(ctx.Html)
	if len(description) > 0 {
		tdki.Description = description[0]
	}
	return tdki
}

func (w *websiteDao) CollectDNS() {

}

func (w *websiteDao) CollectIPAddr() {

}

func (w *websiteDao) CollectSSLCertificateInfo() {

}

func (w *websiteDao) CollectWhois() {

}

func (w *websiteDao) AlertList() {

}

func (w *websiteDao) GetWebSiteUrl() {

}

func (w *websiteDao) SetUrlPoint() {

}

func (w *websiteDao) GetUrlPoint() {

}
