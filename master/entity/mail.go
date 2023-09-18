package entity

// Mail 邮件
type Mail struct {
	From     string `json:"from"`     // 发件人
	AuthCode string `json:"authCode"` // 发件人认证
	Host     string `json:"host"`     // 邮件服务 例子:smtp.qq.com
	Port     int    `json:"port"`     // 邮件服务端口 例子:465/587/25
	ToList   string `json:"toList"`   // 收件人,多个分号";"逗号隔开
}
