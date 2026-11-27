package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	"website-monitor/master/constname"
	"website-monitor/master/entity"
)

type AlertEr interface {
	Get(alertId string) (*entity.AlertData, error)                   // 指定获取报警信息
	GetList() ([]*entity.AlertData, error)                           // 获取所有报警信息列表
	Page(page, size int, filter AlertFilter) (*AlertPageData, error) // 分页查询报警信息
	GetAtWebsite(hostId string) ([]*entity.AlertData, error)         // 获取指定网站的报警信息
	GetWebsiteAlertIDList(hostId string) ([]string, error)           // 获取指定网站报警信息的id
	GetAllWebsiteAlertCounts() (map[string]int, error)               // 一次读取所有站点的报警数量
	Read(id string) error                                            // 消息标记为已读
	ReadAll() error                                                  // 全部标记为已读
	Del(id string) error                                             // 删除报警信息
	Clear(hostId string) error                                       // 清空网站的报警信息
	ClearAll() error                                                 // 清空所有报警信息
}

// AlertFilter 报警列表筛选条件
type AlertFilter struct {
	Seg  string // all / unread / timeout / code
	Host string // hostId
	Kw   string // 站点、URI、报警内容
}

// AlertHost 报警涉及的站点
type AlertHost struct {
	ID   string `json:"id"`
	Host string `json:"host"`
}

// AlertPageData 一页报警数据 + 全局统计信息
type AlertPageData struct {
	List   []*entity.AlertData
	Total  int
	Unread int
	Counts map[string]int
	Hosts  []AlertHost
}

func NewAlert() AlertEr {
	return new(alertDao)
}

type alertDao struct {
}

func (a *alertDao) set(data *entity.AlertData) error {

	err := DB.Set(AlertTable, data.AlertId, data)
	if err != nil {
		return err
	}

	list, err := a.getAtHostID(data.HostId)
	if err != nil && !errors.Is(err, ISNULL) {
		log.Error("getAtHostID err = ", err)
		return err
	}

	list = append(list, data.AlertId)

	err = DB.Set(AlertWebsiteTable, data.HostId, list)
	if err != nil {
		return err
	}

	return nil
}

func (a *alertDao) update(data *entity.AlertData) error {
	return DB.Set(AlertTable, data.AlertId, data)
}

func (a *alertDao) getAtHostID(hostId string) ([]string, error) {
	data := make([]string, 0)

	err := DB.Get(AlertWebsiteTable, hostId, &data)
	if err != nil {
		log.Error("getAtHostID err  = ", err)
	}

	return data, err
}

func (a *alertDao) Get(alertId string) (*entity.AlertData, error) {
	data := &entity.AlertData{}
	err := DB.Get(AlertTable, alertId, &data)
	if err == nil {
		data.AlertType = alertTypeOf(data)
	}
	return data, err
}

func (a *alertDao) GetList() ([]*entity.AlertData, error) {
	list := make([]*entity.AlertData, 0)

	err := DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(AlertTable))
		if b == nil {
			return ISNULL
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			data := &entity.AlertData{}
			err := json.Unmarshal(v, &data)

			if err != nil {
				return err
			}
			data.AlertType = alertTypeOf(data)
			list = append(list, data)

		}

		return nil
	})

	return list, err
}

func (a *alertDao) GetAtWebsite(hostId string) ([]*entity.AlertData, error) {
	data := make([]*entity.AlertData, 0)

	alertIDList, err := a.getAtHostID(hostId)
	if err != nil {
		return data, err
	}

	alertList, err := a.GetList()
	if err != nil {
		return data, err
	}

	for _, id := range alertIDList {

		for _, v := range alertList {
			if id == utils.AnyToString(v.AlertId) {
				data = append(data, v)
			}
		}

	}

	return data, nil
}

func (a *alertDao) GetWebsiteAlertIDList(hostId string) ([]string, error) {
	return a.getAtHostID(hostId)
}

// GetAllWebsiteAlertCounts 一次读取 alert_website_table，避免站点列表 N+1 次 DB 读
func (a *alertDao) GetAllWebsiteAlertCounts() (map[string]int, error) {
	counts := make(map[string]int)

	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AlertWebsiteTable))
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			ids := make([]string, 0)
			if err := json.Unmarshal(v, &ids); err != nil {
				return nil // 跳过脏数据，不影响其他站点
			}
			counts[string(k)] = len(ids)
			return nil
		})
	})

	return counts, err
}

// alertTypeOf 兼容老数据：AlertType 字段是后加的，老记录为空，按消息内容推断。
func alertTypeOf(v *entity.AlertData) string {
	if v.AlertType != "" {
		return v.AlertType
	}
	switch {
	case strings.Contains(v.Msg, "超时"):
		return constname.AlertTypeTimeout
	case strings.Contains(v.Msg, "响应时间超过"):
		// 慢响应在探针侧也归为 timeout 类型
		return constname.AlertTypeTimeout
	case strings.Contains(v.Msg, "状态码"):
		return constname.AlertTypeCode
	case strings.Contains(v.Msg, "请求失败"):
		return constname.AlertTypeErr
	case strings.Contains(v.Msg, "证书"):
		return "ssl"
	}
	return ""
}

// Page 分页查询报警。
// 报警总量在当前量级下很小（千级），一次性读出后在内存中过滤/排序，
// 保证筛选、计数、站点下拉都是全量准确；对外只吐 size 条。
func (a *alertDao) Page(page, size int, filter AlertFilter) (*AlertPageData, error) {
	all, err := a.GetList()
	if err != nil {
		return nil, err
	}

	out := &AlertPageData{
		List:   make([]*entity.AlertData, 0),
		Counts: map[string]int{"all": len(all), "unread": 0, "timeout": 0, "code": 0},
		Hosts:  make([]AlertHost, 0),
	}
	hostMap := make(map[string]string, len(all))
	for _, v := range all {
		if v.Read == 0 {
			out.Unread++
			out.Counts["unread"]++
		}
		switch v.AlertType {
		case constname.AlertTypeTimeout:
			out.Counts["timeout"]++
		case constname.AlertTypeCode:
			out.Counts["code"]++
		}
		hostMap[v.HostId] = v.Host
	}
	for id, host := range hostMap {
		out.Hosts = append(out.Hosts, AlertHost{ID: id, Host: host})
	}
	sort.Slice(out.Hosts, func(i, j int) bool { return out.Hosts[i].Host < out.Hosts[j].Host })

	kw := strings.ToLower(strings.TrimSpace(filter.Kw))
	list := make([]*entity.AlertData, 0, len(all))
	for _, v := range all {
		switch filter.Seg {
		case "unread":
			if v.Read != 0 {
				continue
			}
		case "timeout":
			if v.AlertType != constname.AlertTypeTimeout {
				continue
			}
		case "code":
			if v.AlertType != constname.AlertTypeCode {
				continue
			}
		}
		if filter.Host != "" && v.HostId != filter.Host {
			continue
		}
		if kw != "" {
			text := strings.ToLower(v.Host + " " + v.Uri + " " + v.Msg)
			if !strings.Contains(text, kw) {
				continue
			}
		}
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Date > list[j].Date })

	out.Total = len(list)
	if size <= 0 {
		out.List = list
		return out, nil
	}
	if page < 1 {
		page = 1
	}
	start := (page - 1) * size
	if start >= out.Total {
		return out, nil
	}
	end := start + size
	if end > out.Total {
		end = out.Total
	}
	out.List = list[start:end]
	return out, nil
}

func (a *alertDao) Read(id string) error {
	alert, err := a.Get(id)
	if err != nil {
		return err
	}

	alert.Read = 1 // 0:未读  1:已读

	return a.update(alert)
}

// ReadAll 一次性把所有未读报警标记为已读。
// 放在单个写事务里，避免逐条 DB.Set 触发 N 次 fsync。
func (a *alertDao) ReadAll() error {
	return DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AlertTable))
		if b == nil {
			return fmt.Errorf("%s 表不存在", AlertTable)
		}

		keys := make([][]byte, 0)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			alert := &entity.AlertData{}
			if err := json.Unmarshal(v, alert); err != nil {
				return err
			}
			if alert.Read == 0 {
				keys = append(keys, append([]byte(nil), k...))
			}
		}

		for _, k := range keys {
			raw := b.Get(k)
			if raw == nil {
				continue
			}
			alert := &entity.AlertData{}
			if err := json.Unmarshal(raw, alert); err != nil {
				return err
			}
			alert.Read = 1
			value, err := utils.AnyToJsonB(alert)
			if err != nil {
				return err
			}
			if err := b.Put(k, value); err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *alertDao) Del(id string) error {
	alert, err := a.Get(id)
	if err != nil {
		return err
	}

	websiteAlert, err := a.getAtHostID(alert.HostId)
	if err != nil {
		log.Error("GetAtWebsite err = ", err)
		return err
	}

	for n, v := range websiteAlert {
		if id == v {
			websiteAlert = append(websiteAlert[:n], websiteAlert[n+1:]...)
			break
		}
	}

	err = DB.Set(AlertWebsiteTable, alert.HostId, websiteAlert)
	if err != nil {
		return err
	}

	return DB.Delete(AlertTable, id)
}

func (a *alertDao) Clear(hostId string) error {
	alertList, err := a.getAtHostID(hostId)
	if err != nil {
		return err
	}
	for _, v := range alertList {
		_ = DB.Delete(AlertTable, v)
	}
	return DB.Delete(AlertWebsiteTable, hostId)
}

func (a *alertDao) ClearAll() error {
	return DB.ClearTable(AlertTable)
}

// AlertTimeInfoMap 记录报警超时次数
var AlertTimeInfoMap sync.Map

func getAlertTimeInfoData(key string) []*entity.MonitorLog {
	v, ok := AlertTimeInfoMap.Load(key)
	if ok {
		return v.([]*entity.MonitorLog)
	}

	return make([]*entity.MonitorLog, 0)
}

func setAlertTimeInfoData(key string, data []*entity.MonitorLog) {
	AlertTimeInfoMap.Store(key, data)
}

func AddAlert(mLog *entity.MonitorLog) {
	// 记录报警
	alert := &entity.AlertData{
		AlertId:         utils.IDStr(),
		HostId:          mLog.HostId,
		Host:            mLog.Host,
		AlertType:       mLog.AlertType,
		Date:            mLog.Time,
		Uri:             mLog.Uri,
		UriCode:         mLog.UriCode,
		UriMs:           mLog.UriMs,
		UriType:         mLog.UriType,
		ContrastUri:     mLog.ContrastUri,
		ContrastUriCode: mLog.ContrastUriCode,
		ContrastUriMs:   mLog.ContrastUriMs,
		Ping:            mLog.Ping,
		PingMs:          mLog.PingMs,
		Msg:             mLog.Msg,
		MonitorName:     mLog.MonitorName,
		MonitorIP:       mLog.MonitorIP,
		MonitorAddr:     mLog.MonitorAddr,
		ClientIP:        mLog.ClientIP,
		Read:            0,
	}

	err := new(alertDao).set(alert)
	if err != nil {
		log.Error(err)
	}

	if mLog.AlertType == constname.AlertTypeErr || mLog.AlertType == constname.AlertTypeCode {
		log.Info("发送报警邮件...")
		alertBody := NewAlertBody(mLog)
		NewMail().Send("监测报警通知!", alertBody)

	} else if mLog.AlertType == constname.AlertTypeTimeout {
		log.Error("报警信息为超时...")
		key := mLog.HostId + mLog.MonitorIP // hostID +  MonitorIP
		data := getAlertTimeInfoData(key)

		// 获取报警规则
		rule, err := NewWebsite().GetAlarmRule(mLog.HostId)
		if err != nil {
			log.Error(err)
			return
		}

		// 获取监测设置
		setting, err := NewWebsite().Select(mLog.HostId)
		if err != nil {
			log.Error(err)
			return
		}

		// 判断连续性
		log.Info("data len = ", len(data))
		if len(data) > 0 {

			log.Error("data last = ", data[len(data)-1])
			t1 := -(utils.Date2Timestamp(data[len(data)-1].Time) - time.Now().Unix())
			log.Error("需要执行判断连续性...")
			log.Error("t1 = ", t1)
			t2 := int64(setting.MonitorRate) * 2 // 计算时间段(90%的两个时间段)
			log.Error("t2 = ", t2)

			if t1 > t2 {
				log.Error("超过了连续性...")
				data = make([]*entity.MonitorLog, 0)
			} else {
				data = append(data, mLog)
			}

		} else if len(data) == 0 {
			data = append(data, mLog)
		}

		// 判断连续次数是否触发报警发送邮件
		if len(data) >= int(rule.WebsiteSlowResponseCount) {
			log.Info("发送报警邮件...")
			if alertBody, err := NewAlertBodyAtList(data); err == nil {
				NewMail().Send("监测报警通知!", alertBody)
			}
			data = make([]*entity.MonitorLog, 0)
		}

		setAlertTimeInfoData(key, data)

	}

}

// AlertBody 报警通知
type AlertBody struct {
	Synopsis string
	Tr       []*AlertTd
}

type AlertTd struct {
	Date       string
	Host       string
	Uri        string
	Code       int
	Ms         int64
	NetworkEnv string
	Msg        string
	Monitor    string
}

func (a *AlertBody) Html() string {
	body := ""
	synopsis := fmt.Sprintf("<h3>%s</h3>", a.Synopsis)
	tr := ""
	for _, v := range a.Tr {
		tr += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%d</td><td>%dms</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			v.Date, v.Host, v.Uri, v.Code, v.Ms, v.NetworkEnv, v.Msg, v.Monitor)
	}

	thead := `<thead><tr>
		<th width="auto">监测时间</th>
		<th width="auto">站点</th>
		<th width="auto">链接</th>
		<th width="auto">请求状态码</th>
		<th width="auto">响应时间</th>
		<th width="auto">网络环境</th>
		<th width="auto">报警信息</th>
		<th width="auto">监测器</th>
	</tr></thead>`

	table := fmt.Sprintf(`<table border="1" cellspacing="0">%s<tbody>%s</tbody></table>`, thead, tr)
	body = synopsis + table
	return body
}

func createAlertBody(mLog *entity.MonitorLog) *AlertTd {
	td := &AlertTd{
		Date: mLog.Time,
		Host: mLog.Host,
		Uri:  mLog.Uri,
		Code: mLog.UriCode,
		Ms:   mLog.UriMs,
		Msg:  mLog.Msg,
	}

	td.NetworkEnv = fmt.Sprintf("对照组=> %s ms:%d | Ping=> %s ms:%d ",
		mLog.ContrastUri, mLog.ContrastUriMs, mLog.Ping, mLog.PingMs)
	td.Monitor = fmt.Sprintf("%s| %s| %s", mLog.MonitorName, mLog.MonitorIP, mLog.MonitorAddr)
	return td
}

func NewAlertBody(mLog *entity.MonitorLog) string {
	alert := &AlertBody{
		Synopsis: "监测到" + mLog.Host + "网站出现问题，请快速前往检查并处理!",
		Tr:       make([]*AlertTd, 0),
	}

	alert.Tr = append(alert.Tr, createAlertBody(mLog))

	return alert.Html()
}

func NewAlertBodyAtList(mLogs []*entity.MonitorLog) (string, error) {
	if len(mLogs) < 1 {
		return "", fmt.Errorf("报警日志为空")
	}
	alert := &AlertBody{
		Synopsis: "监测到" + mLogs[0].Host + "等网站出现问题，请快速前往检查并处理!",
		Tr:       make([]*AlertTd, 0),
	}

	for _, v := range mLogs {
		alert.Tr = append(alert.Tr, createAlertBody(v))
	}

	return alert.Html(), nil
}
