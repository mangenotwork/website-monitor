package dao

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	"sync"
	"time"
	"website-monitor/master/constname"
	"website-monitor/master/entity"
)

type AlertEr interface {
	Get(alertId string) (*entity.AlertData, error)           // 指定获取报警信息
	GetList() ([]*entity.AlertData, error)                   // 获取所有报警信息列表
	GetAtWebsite(hostId string) ([]*entity.AlertData, error) // 获取指定网站的报警信息
	GetWebsiteAlertIDList(hostId string) ([]string, error)   // 获取指定网站报警信息的id
	Read(id string) error                                    // 消息标记为已读
	Del(id string) error                                     // 删除报警信息
	Clear(hostId string) error                               // 清空网站的报警信息
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
	if err != nil && err != ISNULL {
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
	return data, err
}

func (a *alertDao) GetList() ([]*entity.AlertData, error) {
	conn := GetDBConn()
	defer func() {
		_ = conn.Close()
	}()
	list := make([]*entity.AlertData, 0)
	err := conn.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(AlertTable))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			data := &entity.AlertData{}
			//log.Info("k = ", string(k))
			err := json.Unmarshal(v, &data)
			if err != nil {
				return err
			}
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
	log.Error("alertIDList = ", alertIDList)
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

func (a *alertDao) Read(id string) error {
	log.Error("read = ", id)
	alert, err := a.Get(id)
	if err != nil {
		return err
	}
	alert.Read = 1 // 0:未读  1:已读
	return a.update(alert)
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
