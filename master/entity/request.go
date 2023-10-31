package entity

// RequestTool 请求器
type RequestTool struct {
	ID                 string            `json:"id"`                 // id
	Name               string            `json:"name"`               // 请求名称
	Notes              string            `json:"notes"`              // 请求描述
	Url                string            `json:"url"`                //
	Method             string            `json:"method"`             // 请求方法
	UserAgent          string            `json:"userAgent"`          // 指定userAgent
	UserAgentRandom    int               `json:"userAgentRandom"`    // userAgent 随机
	Header             map[string]string `json:"header"`             // header
	Query              map[string]string `json:"query"`              // url上的请求参数
	BodyFormData       map[string]string `json:"bodyFormData"`       // multipart/from-data
	BodyFromUrlEncoded map[string]string `json:"bodyFromUrlEncoded"` // application/x-www-from-urlencoded
	BodyJson           string            `json:"bodyJson"`           // application/json
	BodyText           string            `json:"bodyText"`           // text/plain
	IsOpenRetry        int               `json:"isOpenRetry"`        // 是否开启重试 0:关  1:开
	RetryItems         int               `json:"retryItems"`         // 重试次数
	RequestTime        int64             `json:"requestTime"`        // 请求时间 什么时候开始请求的
	ResponseCode       int               `json:"responseCode"`       // Response code
	ResponseTime       int               `json:"responseTime"`       // Response time
	ResponseHeader     string            `json:"responseHeader"`     // 响应头
	ResponseCookie     string            `json:"responseCookie"`     //  Response cookie
	ResponseBody       string            `json:"responseBody"`       // Response body text
}

// RequestToolDir 请求目录
type RequestToolDir struct {
	Name         string            `json:"name"`         // 目录名
	Notes        string            `json:"notes"`        // 目录描述
	PublicHeader map[string]string `json:"publicHeader"` // 公共Header
	Created      string            `json:"created"`      // Created
}

// CookieManage cookie管理
type CookieManage struct {
	Host   string            `json:"host"`   // host
	Cookie map[string]string `json:"cookie"` // cookie
}
