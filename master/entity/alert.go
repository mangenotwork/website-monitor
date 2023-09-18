package entity

// AlertData 报警信息
type AlertData struct {
	Date          string // 监测的时间
	Uri           string // 出现问题的URI
	UriCode       int    // URI响应码
	UriMs         int64  // URI响应时间
	ContrastUriMs int64  // 对照组URI响应时间
	PingMs        int64  // ping 当前网络情况
	Msg           string // 报警信息
}
