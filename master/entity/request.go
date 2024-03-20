package entity

import "net/http"

// RequestTool 请求器
type RequestTool struct {
	ID                 string              `json:"id"`                 // id
	Name               string              `json:"name"`               // 请求名称
	Note               string              `json:"note"`               // 请求描述
	Url                string              `json:"url"`                //
	Method             string              `json:"method"`             // 请求方法
	UserAgent          string              `json:"userAgent"`          // 指定userAgent
	UserAgentRandom    int                 `json:"userAgentRandom"`    // userAgent 随机
	Header             map[string]string   `json:"header"`             // header
	Query              map[string]string   `json:"query"`              // url上的请求参数
	BodyFormData       map[string]string   `json:"bodyFormData"`       // multipart/from-data
	BodyFromUrlEncoded map[string]string   `json:"bodyFromUrlEncoded"` // application/x-www-from-urlencoded
	BodyJson           string              `json:"bodyJson"`           // application/json
	BodyText           string              `json:"bodyText"`           // text/plain
	IsOpenRetry        int                 `json:"isOpenRetry"`        // 是否开启重试 0:关  1:开
	RetryItems         int                 `json:"retryItems"`         // 重试次数
	RequestTime        string              `json:"reqTime"`            // 请求时间 什么时候开始请求的
	ResponseCode       int                 `json:"respCode"`           // Response code
	ResponseTime       string              `json:"respTime"`           // Response time
	ResponseHeader     map[string][]string `json:"respHeader"`         // 响应头
	ResponseCookie     []*http.Cookie      `json:"respCookie"`         //  Response cookie
	ResponseBody       string              `json:"respBody"`           // Response body text
	RequestHeader      map[string][]string `json:"reqHeader"`          // request header
	HostIP             string              `json:"hostIP"`             // 请求主机的ip
	HostIPAddr         string              `json:"hostIPAddr"`         // 请求主机的ip属地
	ClientName         string              `json:"clientName"`         // 请求器昵称
	ClientIP           string              `json:"clientIP"`           // 内网
	ClientPublicIP     string              `json:"clientPublicIP"`     // 公网
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

// RequesterGlobalHeader 全局header
type RequesterGlobalHeader struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Enable bool   `json:"enable"`
	Note   string `json:"note"` // 描述
}

// RequestNowList 当前请求列表
type RequestNowList struct {
	Id     string `json:"id"`
	Method string `json:"method"`
	Url    string `json:"url"`
	Name   string `json:"name"`
	IsNow  bool   `json:"isNow"`
	Time   int64  `json:"time"` // 主要用于排序
}
