package business

import (
	"bytes"
	"fmt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	gt "github.com/mangenotwork/gathertool"
	udp "github.com/mangenotwork/udp_comm"
	"sync"
	"time"
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
	Host                    string `json:"host"`
	HostID                  string `json:"hostID"`
	MonitorRate             int64  `json:"monitorRate"`
	ContrastUrl             string `json:"contrastUrl"`
	ContrastTime            int64  `json:"contrastTime"`
	Ping                    string `json:"ping"`
	PingTime                int64  `json:"pingTime"`
	Notes                   string `json:"notes"`
	Created                 int64  `json:"created"`
	WebsiteSlowResponseTime int    `json:"websiteSlowResponseTime"`
}

func (item *WebsiteItem) Add() {
	AllWebsiteData.Store(item.HostID, item)
}

func (item *WebsiteItem) Put(data []byte) {
	item.Conn.Put("monitor", data)
}

func (item *WebsiteItem) mLogSerialize(mLog *MonitorLog) []byte {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%d|%d|%s|%d|%d|%s|%d|%s|%s|%s|%s|%s|",
		mLog.LogType, mLog.Time, mLog.HostId, mLog.Host, mLog.UriType, mLog.Uri, mLog.UriCode, mLog.UriMs,
		mLog.ContrastUri, mLog.ContrastUriCode, mLog.ContrastUriMs, mLog.Ping, mLog.PingMs, mLog.Msg, mLog.AlertType,
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
func (item *WebsiteItem) AlertRuleMs(resMs int64) bool {
	if resMs > int64(item.WebsiteSlowResponseTime) {
		return true
	}
	return false
}

func (item *WebsiteItem) MonitorHealthUri(mLog *MonitorLog) {
	mLog.Uri = item.Host
	mLog.UriType = URIHealth
	mLog.LogType = LogTypeInfo
	mLog.Msg = ""
	mLog.AlertType = AlertTypeNone
	log.Info("=================================  监测生命URI... ", item.Host)
	healthCode, healthMs, err := request(item.Host)
	if err != nil {
		mLog.LogType = LogTypeAlert
		mLog.Msg = "请求失败，err=" + err.Error()
		mLog.AlertType = AlertTypeErr
		item.Put(item.mLogSerialize(mLog))
		return
	}
	mLog.UriCode = healthCode
	mLog.UriMs = healthMs
	// 监测规则
	if item.AlertRuleCode(healthCode) {
		mLog.LogType = LogTypeAlert
		mLog.Msg = fmt.Sprintf("请求失败，状态码:%d;", healthCode)
		mLog.AlertType = AlertTypeCode
		item.Put(item.mLogSerialize(mLog))
		return
	}
	// 请求超时报警
	if item.AlertRuleMs(healthMs) {
		mLog.LogType = LogTypeAlert
		mLog.Msg += fmt.Sprintf("响应时间超过设置的报警时间，响应时间:%dms;", healthMs)
		mLog.AlertType = AlertTypeTimeout
		item.Put(item.mLogSerialize(mLog))
		return
	}
	mLog.Msg = "passed"
	item.Put(item.mLogSerialize(mLog))
	return
}

func (item *WebsiteItem) MonitorRandomUri(mLog *MonitorLog) {
	uri := GetWebsiteUrlDataMap(item.HostID)
	if len(uri) > 0 {
		time.Sleep(1 * time.Second) // 强行休息1s
		mLog.UriType = URIRandom
		mLog.LogType = LogTypeInfo     // 复位
		mLog.Msg = ""                  // 复位
		mLog.AlertType = AlertTypeNone // 复位
		randomUri := utils.RandomString(uri)
		mLog.Uri = randomUri
		log.Info("=================================  随机取一个URI监测... ", randomUri)
		randomCode, randomMs, err := request(randomUri)
		if err != nil {
			mLog.LogType = LogTypeAlert
			mLog.Msg = "请求失败，err=" + err.Error()
			mLog.AlertType = AlertTypeErr
			item.Put(item.mLogSerialize(mLog))
			return
		}
		mLog.UriCode = randomCode
		mLog.UriMs = randomMs
		if item.AlertRuleCode(randomCode) {
			mLog.LogType = LogTypeAlert
			mLog.Msg = fmt.Sprintf("请求失败，状态码:%d", randomCode)
			mLog.AlertType = AlertTypeCode
			item.Put(item.mLogSerialize(mLog))
			return
		}
		if item.AlertRuleMs(randomMs) {
			mLog.LogType = LogTypeAlert
			mLog.Msg += fmt.Sprintf("响应时间超过设置的报警时间，响应时间:%dms", randomMs)
			mLog.AlertType = AlertTypeTimeout
			item.Put(item.mLogSerialize(mLog))
			return
		}
		mLog.Msg = "passed"
		item.Put(item.mLogSerialize(mLog))
	}
}

func (item *WebsiteItem) MonitorPointUri(mLog *MonitorLog, pointUrl string) {
	log.Info("=================================  执行监测点监测... ", pointUrl)
	time.Sleep(1 * time.Second) // 强行休息1s
	mLog.UriType = URIPoint
	mLog.LogType = LogTypeInfo     // 复位
	mLog.Msg = ""                  // 复位
	mLog.AlertType = AlertTypeNone // 复位
	mLog.Uri = pointUrl
	pointCode, pointMs, err := request(pointUrl)
	if err != nil {
		mLog.LogType = LogTypeAlert
		mLog.Msg = "请求失败，err=" + err.Error()
		mLog.AlertType = AlertTypeErr
		item.Put(item.mLogSerialize(mLog))
		return
	}
	mLog.UriCode = pointCode
	mLog.UriMs = pointMs
	// 监测规则
	if item.AlertRuleCode(pointCode) {
		mLog.LogType = LogTypeAlert
		mLog.Msg = fmt.Sprintf("请求失败，状态码:%d", pointCode)
		mLog.AlertType = AlertTypeCode
		item.Put(item.mLogSerialize(mLog))
		return
	}
	if item.AlertRuleMs(pointMs) {
		mLog.LogType = LogTypeAlert
		mLog.Msg += fmt.Sprintf("响应时间超过设置的报警时间，响应时间:%dms", pointMs)
		mLog.AlertType = AlertTypeTimeout
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
