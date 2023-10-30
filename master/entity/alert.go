package entity

// AlertData 报警信息
type AlertData struct {
	AlertId         string `json:"alertId"`
	HostId          string `json:"hostId"`
	Host            string `json:"host"`
	Date            string `json:"date"`            // 监测的时间
	Uri             string `json:"uri"`             // 出现问题的URI
	UriCode         int    `json:"uriCode"`         // URI响应码
	UriMs           int64  `json:"uriMs"`           // URI响应时间
	UriType         string `json:"uriType"`         // 监测的URI类型 Health:根URI,健康URI  Random:随机URI  Point:监测点URI
	ContrastUri     string `json:"contrastUri"`     // 对照组URI
	ContrastUriCode int    `json:"contrastUriCode"` // 对照组URI响应码
	ContrastUriMs   int64  `json:"contrastUriMs"`   // 对照组URI响应时间
	Ping            string `json:"ping"`
	PingMs          int64  `json:"pingMs"`
	Msg             string `json:"msg"`         // 报警信息
	MonitorName     string `json:"monitorName"` // 监测器名称
	MonitorIP       string `json:"monitorIP"`   // 监测器 公网ip
	MonitorAddr     string `json:"monitorAddr"` // 监测器 ip属地
	ClientIP        string `json:"clientIP"`    // 监测器连接地址
	Read            int    `json:"read"`        // 报警信息已读  0: 未读  1: 已读
}
