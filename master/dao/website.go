package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"website-monitor/master/entity"

	"github.com/boltdb/bolt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	gt "github.com/mangenotwork/gathertool"
)

type WebsiteEr interface {
	Add(data *entity.Website, alarmRule *entity.WebsiteAlarmRule, scan *entity.WebsiteScanCheckUp) error

	Del(hostID string) error

	Edit(base *entity.Website, alarmRule *entity.WebsiteAlarmRule, scan *entity.WebsiteScanCheckUp) error

	SelectList() ([]*entity.Website, int, error)

	Select(hostID string) (*entity.Website, error)

	// GetConfAlarmRule 获取网站的监测配置
	GetConfAlarmRule(hostID string) (*entity.WebsiteAlarmRule, error)

	GetConfScanCheckUp(hostID string) (*entity.WebsiteScanCheckUp, error)

	// GetAlarmRule 获取监测报警规则
	GetAlarmRule(hostID string) (*entity.WebsiteAlarmRule, error)

	// GetScanCheckUp 获取扫描规则
	GetScanCheckUp(hostID string) (*entity.WebsiteScanCheckUp, error)

	// AlertList 报警列表
	AlertList()

	// GetWebSiteUrl 获取采集的网站Url
	GetWebSiteUrl(hostId string) (*entity.WebSiteUrl, error)

	// Collect 采集网站信息
	Collect(host string) *entity.WebsiteInfo

	// SaveCollectInfo 保存采集网站信息
	SaveCollectInfo(host, hostID string) error

	// GetInfo 获取网站信息
	GetInfo(hostID string) (*entity.WebsiteInfo, error)

	// CollectTDK 采集网站页面基础信息 - 刷新功能
	CollectTDK(host string) *entity.TDKI

	// CollectDNS 采集网站DNS信息 - 刷新功能
	CollectDNS(host string) error

	// CollectIPAddr 采集网站IP和属地 - 刷新功能
	CollectIPAddr(host string) error

	// CollectSSLCertificateInfo 采集网站证书信息 - 刷新功能
	CollectSSLCertificateInfo(host string) error

	// CollectWhois 采集网站Whois信息 - 刷新功能
	CollectWhois(host string) error

	// CollectICP 采集网站ipc信息 - 刷新功能
	CollectICP(host string) error

	// SetPoint  设置网站监测点url
	SetPoint(hostID, url string) error

	// GetPoint  获取网站监测点url
	GetPoint(hostID string) (*entity.WebSitePoint, error)

	// DelPoint  删除指定网站监测点url
	DelPoint(hostID, url string) error

	// ClearPoint 清空指定网站监测点url
	ClearPoint(hostID string) error
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

	hostKey := utils.GetMD5Encode(data.Host)
	data.HostID = hostKey

	// 判断是否存在
	has, _ := w.Select(hostKey)
	if has.Host == data.Host {
		return fmt.Errorf("网站:%s 已创建监测, 无法重复添加。", data.Host)
	}

	err = w.addWebsite(data)
	if err != nil {
		return err
	}

	// 保存报警规则信息
	rule.Host = data.Host
	rule.HostID = hostKey

	err = w.addWebsiteAlarmRule(rule)
	if err != nil {
		return err
	}

	// 保存扫描规则信息
	scan.Host = data.Host
	scan.HostID = hostKey
	err = w.addWebsiteScanCheckUp(scan)
	if err != nil {
		return err
	}

	// 异步获取网站信息
	go func() {
		_ = w.SaveCollectInfo(data.Host, hostKey)
	}()

	// 异步执行扫描
	go func() {
		Scan(data.Host, hostKey, scan.ScanDepth)
	}()

	// 异步通知监测器更新监测网站
	go func() {
		NoticeUpdateWebsite(data.HostID)
	}()

	return err
}

func (w *websiteDao) GetAlarmRule(hostID string) (*entity.WebsiteAlarmRule, error) {
	value := &entity.WebsiteAlarmRule{}
	err := DB.Get(WebsiteAlarmRuleTable, hostID, &value)
	return value, err
}

func (w *websiteDao) GetScanCheckUp(hostID string) (*entity.WebsiteScanCheckUp, error) {
	value := &entity.WebsiteScanCheckUp{}
	err := DB.Get(WebsiteScanCheckUpTable, hostID, &value)
	return value, err
}

func (w *websiteDao) addWebsite(data *entity.Website) error {
	return DB.Set(WebSiteTable, data.HostID, data)
}

func (w *websiteDao) addWebsiteAlarmRule(alarmRule *entity.WebsiteAlarmRule) error {
	if alarmRule.Host == "" {
		return fmt.Errorf("没有host")
	}
	return DB.Set(WebsiteAlarmRuleTable, alarmRule.HostID, alarmRule)
}

func (w *websiteDao) addWebsiteScanCheckUp(scan *entity.WebsiteScanCheckUp) error {
	if scan.Host == "" {
		return fmt.Errorf("没有host")
	}
	return DB.Set(WebsiteScanCheckUpTable, scan.HostID, scan)
}

func (w *websiteDao) Del(hostID string) error {
	// 删除 Website
	err := DB.Delete(WebSiteTable, hostID)

	// 删除 WebsiteAlarmRule
	err = DB.Delete(WebsiteAlarmRuleTable, hostID)

	// 删除 WebsiteScanCheckUp
	err = DB.Delete(WebsiteScanCheckUpTable, hostID)

	// 删除 WebsiteInfo
	err = DB.Delete(WebSiteInfoTable, hostID)

	// 删除监测点 Website
	err = DB.Delete(WebSiteUrlPointTable, hostID)

	// 通知更新监测网站
	go func() {
		NoticeDelWebsite(hostID)
	}()

	return err
}

func (w *websiteDao) Edit(base *entity.Website, alarmRule *entity.WebsiteAlarmRule, scan *entity.WebsiteScanCheckUp) error {
	err := w.addWebsite(base)
	err = w.addWebsiteAlarmRule(alarmRule)
	err = w.addWebsiteScanCheckUp(scan)
	return err
}

func (w *websiteDao) SelectList() ([]*entity.Website, int, error) {
	DB.Open()

	defer func() {
		_ = DB.Conn.Close()
	}()

	count := 0
	data := make([]*entity.Website, 0)

	err := DB.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(WebSiteTable))
		if b == nil {
			return fmt.Errorf(WebSiteTable + "表不存在")
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			count++
			value := &entity.Website{}

			e := json.Unmarshal(v, value)
			if e != nil {
				log.Error("数据解析错误")
			}

			data = append(data, value)
		}
		return nil
	})

	return data, count, err
}

func (w *websiteDao) Select(host string) (*entity.Website, error) {
	website := &entity.Website{}
	err := DB.Get(WebSiteTable, host, &website)
	return website, err
}

func (w *websiteDao) Collect(host string) *entity.WebsiteInfo {
	log.Info("Collect --> ", host)
	info := &entity.WebsiteInfo{
		Host: host,
	}

	tdkI := w.CollectTDK(host)
	info.Title = tdkI.Title
	info.Keywords = tdkI.Keywords
	info.Description = tdkI.Description
	info.Icon = tdkI.Icon

	// dns
	info.DNS = NsLookUpLocal(host)

	// IPAddr
	info.IPAddr = make([]*entity.IPAddr, 0)
	for _, v := range info.DNS.IPs {
		addr := GetIP(v)
		info.IPAddr = append(info.IPAddr, &entity.IPAddr{IP: v, Address: addr})
	}

	// SSLCertificateInfo
	info.SSLCertificateInfo, _ = GetCertificateInfo(host)

	// Whois
	info.Whois = Whois(host)

	// ipc
	info.IPC = GetICP(host)

	header := w.responseHeaders(host)
	info.Server = header.Get("Server")
	info.ContentEncoding = header.Get("Content-Encoding")
	info.ContentLanguage = header.Get("Content-Language")

	return info
}

// SaveCollectInfo 保存采集的网站信息
func (w *websiteDao) SaveCollectInfo(host, hostID string) error {
	log.Info("data = ", host)
	info := w.Collect(host)
	info.HostID = hostID
	return DB.Set(WebSiteInfoTable, hostID, info)
}

func (w *websiteDao) GetInfo(hostID string) (*entity.WebsiteInfo, error) {
	info := &entity.WebsiteInfo{}
	err := DB.Get(WebSiteInfoTable, hostID, info)
	return info, err
}

func (w *websiteDao) responseHeaders(url string) http.Header {
	ctx, _ := gt.Get(url, gt.Header{
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": "zh-CN,zh;q=0.9",
	})
	return ctx.Resp.Header
}

func (w *websiteDao) CollectTDK(host string) *entity.TDKI {
	host = urlStr(host)
	tdkI := &entity.TDKI{
		Icon: host + "/favicon.ico",
	}

	ctx, _ := gt.Get(host)
	title := gt.RegHtmlTitleTxt(ctx.Html)
	if len(title) > 0 {
		tdkI.Title = title[0]
	}

	keyword := gt.RegHtmlKeywordTxt(ctx.Html)
	if len(keyword) > 0 {
		tdkI.Keywords = keyword[0]
	}

	description := gt.RegHtmlDescriptionTxt(ctx.Html)
	if len(description) > 0 {
		tdkI.Description = description[0]
	}

	return tdkI
}

func (w *websiteDao) CollectDNS(host string) error {
	info, err := w.GetInfo(host)
	if err != nil {
		return err
	}

	info.DNS = NsLookUpLocal(host)
	return DB.Set(WebSiteInfoTable, host, info)
}

func (w *websiteDao) CollectIPAddr(host string) error {
	info, err := w.GetInfo(host)
	if err != nil {
		return err
	}

	ipAddr := make([]*entity.IPAddr, 0)
	for _, v := range info.DNS.IPs {
		addr := GetIP(v)
		ipAddr = append(ipAddr, &entity.IPAddr{IP: v, Address: addr})
	}

	info.IPAddr = ipAddr

	return DB.Set(WebSiteInfoTable, host, info)
}

func (w *websiteDao) CollectSSLCertificateInfo(host string) error {
	info, err := w.GetInfo(host)
	if err != nil {
		return err
	}

	info.SSLCertificateInfo, _ = GetCertificateInfo(host)
	return DB.Set(WebSiteInfoTable, host, info)
}

func (w *websiteDao) CollectWhois(host string) error {
	info, err := w.GetInfo(host)
	if err != nil {
		return err
	}

	info.Whois = Whois(host)
	return DB.Set(WebSiteInfoTable, host, info)
}

func (w *websiteDao) CollectICP(host string) error {
	info, err := w.GetInfo(host)
	if err != nil {
		return err
	}

	info.IPC = GetICP(host)
	return DB.Set(WebSiteInfoTable, host, info)
}

func (w *websiteDao) AlertList() {

}

func (w *websiteDao) GetWebSiteUrl(hostId string) (*entity.WebSiteUrl, error) {
	urlData := &entity.WebSiteUrl{}
	err := DB.Get(WebSiteURITable, hostId, &urlData)
	return urlData, err
}

func (w *websiteDao) SetPoint(hostID, url string) error {
	data, err := w.GetPoint(hostID)

	if err != nil && !errors.Is(err, ISNULL) {
		return err
	}

	for _, v := range data.URL {
		if v == url {
			return fmt.Errorf("监测点存在")
		}
	}

	data.URL = append(data.URL, url)
	data.HostID = hostID

	err = DB.Set(WebSiteUrlPointTable, hostID, data)
	if err == nil {

		go func() {
			NoticeUpdateWebsitePoint(hostID)
		}()
	}

	return err
}

func (w *websiteDao) GetPoint(hostID string) (*entity.WebSitePoint, error) {
	data := &entity.WebSitePoint{
		HostID: hostID,
	}

	err := DB.Get(WebSiteUrlPointTable, hostID, &data)
	return data, err
}

func (w *websiteDao) DelPoint(hostID, url string) error {
	data, err := w.GetPoint(hostID)
	if err != nil && !errors.Is(err, ISNULL) {
		return err
	}

	for n, v := range data.URL {

		if v == url {
			data.URL = append(data.URL[:n], data.URL[n+1:]...)
			break
		}

	}

	err = DB.Set(WebSiteUrlPointTable, hostID, data)
	if err == nil {

		go func() {
			NoticeUpdateWebsitePoint(hostID)
		}()
	}

	return err
}

func (w *websiteDao) ClearPoint(hostID string) error {
	err := DB.Delete(WebSiteUrlPointTable, hostID)
	if err == nil {

		go func() {
			NoticeUpdateWebsitePoint(hostID)
		}()
	}

	return err
}

func (w *websiteDao) GetConfAlarmRule(hostID string) (*entity.WebsiteAlarmRule, error) {
	data := &entity.WebsiteAlarmRule{}
	err := DB.Get(WebsiteAlarmRuleTable, hostID, data)
	return data, err
}

func (w *websiteDao) GetConfScanCheckUp(hostID string) (*entity.WebsiteScanCheckUp, error) {
	data := &entity.WebsiteScanCheckUp{}
	err := DB.Get(WebsiteScanCheckUpTable, hostID, data)
	return data, err
}
