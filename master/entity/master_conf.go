package entity

var (
	ContrastUriDefault  = "www.baidu.com"
	ContrastTimeDefault = 1000
	PingDefault         = "101.226.4.6"
	LogSaveDayDefault   = 7
)

// MasterConf master设置
type MasterConf struct {
	// 监测日志保留天数
	LogSaveDay int
}
