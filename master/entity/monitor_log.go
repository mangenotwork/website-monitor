package entity

// MonitorLog 监测日志
type MonitorLog struct {
	LogType         string // Info  Alert Error
	Time            string
	HostId          string
	Host            string
	UriType         string // 监测的URI类型 Health:根URI,健康URI  Random:随机URI  Point:监测点URI
	Uri             string // URI
	UriCode         int    // URI响应码
	UriMs           int64  // URI响应时间
	ContrastUri     string // 对照组URI
	ContrastUriCode int    // 对照组URI响应码
	ContrastUriMs   int64  // 对照组URI响应时间
	Ping            string
	PingMs          int64
	Msg             string
	AlertType       string
	MonitorName     string // 监测器名称
	MonitorIP       string // 监测器 ip地域信息
	MonitorAddr     string
	ClientIP        string
}
