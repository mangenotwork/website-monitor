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
