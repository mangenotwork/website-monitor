package entity

// AlertData 报警信息
type AlertData struct {
	Date          string `json:"date"`          // 监测的时间
	Uri           string `json:"uri"`           // 出现问题的URI
	UriCode       int    `json:"uriCode"`       // URI响应码
	UriMs         int64  `json:"uriMs"`         // URI响应时间
	ContrastUriMs int64  `json:"contrastUriMs"` // 对照组URI响应时间
	PingMs        int64  `json:"pingMs"`        // ping 当前网络情况
	Msg           string `json:"msg"`           // 报警信息
}
