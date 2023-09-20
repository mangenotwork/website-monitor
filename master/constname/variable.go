package constname

import "time"

var (
	UserToken    string = "sign"
	TokenExpires        = 60 * 60 * 24 * 7
	LastSendMail int64  = 0
	TimeStamp           = time.Now().Unix()
)

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
