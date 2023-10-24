package dao

import (
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	"sync"
	"time"
	"website-monitor/master/constname"
	"website-monitor/master/entity"
)

type AlertEr interface {
	Get()     // 指定获取网站监测的报警信息
	GetList() // 获取报警信息列表
}

type alertDao struct {
}

func (a *alertDao) set(data *entity.AlertData) error {
	return DB.Set(AlertTable, data.HostId, data)
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
		// TODO...
		log.Info("发送报警邮件...")
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
			// TODO...
			log.Info("发送报警邮件...")
			data = make([]*entity.MonitorLog, 0)
		}

		setAlertTimeInfoData(key, data)

	}

}
