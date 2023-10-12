package business

import (
	"bytes"
	"fmt"
	"github.com/mangenotwork/beacon-tower/udp"
	"github.com/mangenotwork/common/log"
	gt "github.com/mangenotwork/gathertool"
	"sync"
)

var AllWebsiteData sync.Map

func EmptyAllWebsiteData() {
	AllWebsiteData = sync.Map{}
}

type WebsiteItem struct {
	*Website
	RateItem  int64 // 用于计算
	Conn      *udp.Client
	LoopPoint int64
}

// Website 监测网站
type Website struct {
	Host         string `json:"host"`
	HostID       string `json:"hostID"`
	MonitorRate  int64  `json:"monitorRate"`
	ContrastUrl  string `json:"contrastUrl"`
	ContrastTime int64  `json:"contrastTime"`
	Ping         string `json:"ping"`
	PingTime     int64  `json:"pingTime"`
	Notes        string `json:"notes"`
	Created      int64  `json:"created"`
}

func (item *WebsiteItem) Add() {
	AllWebsiteData.Store(item.HostID, item)
}

func (item *WebsiteItem) Put(data []byte) {
	item.Conn.Put("monitor", data)
}

func (item *WebsiteItem) mLogSerialize(mLog *MonitorLog) []byte {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%d|%d|%s|%d|%d|%s|%d|%s|%s|%s|%s|",
		mLog.LogType, mLog.Time, mLog.HostId, mLog.Host, mLog.UriType, mLog.Uri, mLog.UriCode, mLog.UriMs,
		mLog.ContrastUri, mLog.ContrastUriCode, mLog.ContrastUriMs, mLog.Ping, mLog.PingMs, mLog.Msg,
		mLog.MonitorName, mLog.MonitorIP, mLog.MonitorAddr))
	return buf.Bytes()
}

func (item *WebsiteItem) PingActive(mLog *MonitorLog) (int64, bool) {
	ping, err := gt.Ping(item.Ping)
	if err != nil {
		mLog.LogType = LogTypeError
		mLog.Msg = "网络不通请前往检查监测平台!" + err.Error()
		item.Put(item.mLogSerialize(mLog))
		return 0, false
	}
	pingMs := ping.Milliseconds()
	mLog.PingMs = pingMs
	if pingMs >= 1000 {
		mLog.LogType = LogTypeError
		mLog.Msg = fmt.Sprintf("网络环境缓慢，超过1s(%d)请前往检查监测平台!", pingMs)
		item.Put(item.mLogSerialize(mLog))
		return pingMs, false
	}
	return pingMs, true
}

func (item *WebsiteItem) ContrastActive(mLog *MonitorLog) bool {
	contrastErr := false
	contrastCode, contrastMs, err := request(item.ContrastUrl)
	if err != nil {
		contrastErr = true
		mLog.Msg += fmt.Sprintf("对照组请求失败 err =%v!", err)
		item.Put(item.mLogSerialize(mLog))
	}
	mLog.ContrastUriCode = contrastCode
	mLog.ContrastUriMs = contrastMs
	if item.AlertRuleCode(contrastCode) {
		contrastErr = true
		mLog.Msg += fmt.Sprintf("对照组请求失败code=%d!", contrastCode)
		item.Put(item.mLogSerialize(mLog))
	}
	if contrastMs >= item.ContrastTime {
		contrastErr = true
		mLog.Msg += fmt.Sprintf("请求对照组网络超%dms 当前为%dms!", item.ContrastTime, contrastMs)
		item.Put(item.mLogSerialize(mLog))
	}
	if contrastErr {
		mLog.LogType = LogTypeError
	}
	return contrastErr
}

// AlertRuleCode 报警规则 出现 400, 404, >500 的状态码视为出现问题
func (item *WebsiteItem) AlertRuleCode(code int) bool {
	if code == 400 || code == 404 || code >= 500 {
		return true
	}
	return false
}

// AlertRuleMs 响应时间超过设置的响应时间视为出现问题
// TODO 获取报警规则
func (item *WebsiteItem) AlertRuleMs(nowMs int64) bool {
	return false
}

func (item *WebsiteItem) MonitorHealthUri(mLog *MonitorLog) {
	mLog.UriType = URIHealth
	mLog.LogType = LogTypeInfo
	mLog.Msg = ""
	// =================================  监测生命URI
	log.Info("=================================  监测生命URI... ", item.Host)
	times := 0
R:
	healthCode, healthMs, err := request(item.Host)
	if err != nil {
		mLog.LogType = LogTypeAlert
		mLog.Msg = "请求失败，err=" + err.Error()
		item.Put(item.mLogSerialize(mLog))
		return
	}
	mLog.Uri = item.Host
	mLog.UriCode = healthCode
	mLog.UriMs = healthMs
	// 监测规则
	if item.AlertRuleCode(healthCode) {
		mLog.LogType = LogTypeAlert
		mLog.Msg = fmt.Sprintf("请求失败，状态码:%d;", healthCode)
		item.Put(item.mLogSerialize(mLog))
		return
	}
	// TODO 请求超时报警业务
	if item.AlertRuleMs(healthMs) {
		// 如果是超时再来一次,确保监测是连续超时并非单个网络波动
		if times == 0 {
			times++
			goto R
		}
		mLog.LogType = LogTypeAlert
		mLog.Msg += fmt.Sprintf("响应时间超过设置的报警时间，响应时间:%dms;", healthMs)
		item.Put(item.mLogSerialize(mLog))
		return
	}
	mLog.Msg = "passed"
	item.Put(item.mLogSerialize(mLog))
	return
}

var WebsiteAllUrlData sync.Map

// WebsiteAllUrl 监测网站的所有Url
type WebsiteAllUrl struct {
	HostID string   `json:"hostID"`
	Url    []string `json:"url"`
}

func (w *WebsiteAllUrl) Add() {
	WebsiteAllUrlData.Store(w.HostID, w)
}
