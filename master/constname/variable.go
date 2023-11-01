package constname

import (
	"time"
)

// sys
var (
	UserToken    string = "sign"
	TokenExpires        = 60 * 60 * 24 * 7
	LastSendMail int64  = 0
	TimeStamp           = time.Now().Unix()
)

// sys
const (
	DayLayout     = "20060102"
	MasterVersion = "v0.1"
)

// default param
var (
	DefaultMonitorRate              int64 = 15
	DefaultContrastUrl                    = "https://www.baidu.com"
	DefaultContrastTime             int64 = 1000
	DefaultPing                           = "8.8.8.8"
	DefaultPingTime                 int64 = 1000
	DefaultUriDepth                 int64 = 2
	DefaultScanRate                 int64 = 24
	DefaultWebsiteSlowResponseCount int64 = 3
	DefaultSSLCertificateExpire     int64 = 14
)

// notice
const (
	NoticeUpdateWebsiteAllLabel    = "websiteAll" // 通知更新网站监测所有
	NoticeUpdateWebsiteLabel       = "website"    // 通知更新网站监测
	NoticeDelWebsiteLabel          = "websiteDel" // 通知删除网站监测
	NoticeUpdateWebsiteAllUrlLabel = "allUrl"     // 通知更新网站url
	NoticeUpdateWebsitePointLabel  = "point"      // 通知更新网站监测点
)

// log & alert
const (
	URIHealth        = "Health"
	URIRandom        = "Random"
	URIPoint         = "Point"
	LogTypeInfo      = "Info"
	LogTypeAlert     = "Alert"
	LogTypeError     = "Error"
	AlertTypeNone    = ""
	AlertTypeErr     = "err"
	AlertTypeCode    = "code"
	AlertTypeTimeout = "timeout"
)

// mail
const (
	MailSendOpen  = 0
	MailSendClose = 1
)
