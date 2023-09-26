package business

import "sync"

var AllWebsite sync.Map

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
