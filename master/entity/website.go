package entity

import (
	"fmt"
	"github.com/mangenotwork/common/log"
	"time"
)

// Website 监测网站
type Website struct {
	// 默认取根url进行检测，主键
	Host string `json:"host"`

	// 监测频率  单位 ms
	MonitorRate int64 `json:"monitorRate"`

	// 保证监测的有效性:每次监测前用过对照组+ping来确认当前监测器网络的情况，只有当网络情况好才执行监测请求
	// 对照组Url
	ContrastUrl string
	// 对照组响应超过这个时间判定为当前网络不稳定不执行本次监测请求
	ContrastTime int64
	// 用于检查当前网络
	Ping string
	// ping响应超过这个时间判定为当前网络不稳定不执行本次监测请求
	PingTime int64

	Notes   string `json:"notes"`
	Created int64  `json:"created"`
}

// WebsiteAlarmRule 报警规则
type WebsiteAlarmRule struct {
	Host                     string // 主键与website Host对应
	WebsiteSlowResponseTime  int64  // 单位ms  网站响应有多慢才记录报警
	WebsiteSlowResponseCount int64  // 连续几次慢就发送邮件通知
	SSLCertificateExpire     int64  // 单位天 证书还有几天过期就触发报警
	NotTDK                   bool   // 在随机监测中发现网站页面没有存在 title, description, keywords 就触发报警
	BadLink                  bool   // 扫描的时候存在死链就报警
	ExtLinkChange            bool   // 扫描间隔判断外链有变化则报警,可以监测到是否被劫持
}

// WebsiteScanCheckUp 网站扫描检查内容
type WebsiteScanCheckUp struct {
	Host         string // 主键与website Host对应
	ScanDepth    int64  `json:"uriDepth"` // 扫描站点深度 默认 2
	ScanRate     int64  `json:"scanRate"` // 扫描网站频率  单位秒
	ScanExtLinks bool   // 是否检查外链,对比上一次扫描数据判别外链是否变化
	ScanBadLink  bool   // 是否扫描死链接
	// TODO... 安全扫描， Sql注入, XSS等等...
}

// WebsiteInfo 网站基本信息
type WebsiteInfo struct {
	Host               string // 主键与website Host对应
	Title              string
	Description        string
	Keywords           string
	Icon               string
	DNS                *DNSInfo            // 网站DNS信息
	IPAddr             []*IPAddr           // 网站IP和属地
	Server             string              // response Headers 读取 headers Server
	ContentEncoding    string              //  response Headers 是否压缩 读取 headers Content-Encoding
	ContentLanguage    string              // response Headers 语言 读取 headers Content-Language
	SSLCertificateInfo *SSLCertificateInfo // 证书信息
	Filing             string              // 网站备案信息
	Whois              *WhoisInfo          // Whois信息
}

// WebSiteUrl 网站的URL存储
type WebSiteUrl struct {
	Host    string   // 主键与website Host对应
	AllUri  []string // 抓取到的所有链接
	ExtLink []string // 外链
	BadLink []string // 死链
	JsLink  []string // 资源文件 js
	CssLink []string // 资源文件 css
}

// WebSiteUrlPoint 指定网站监测Url
type WebSiteUrlPoint struct {
	Host string // 主键与website Host对应
	Uri  []string
}

// WebSiteAlert 监控报警信息
type WebSiteAlert struct {
	Host string // 主键与website Host对应
	List []*AlertData
}

// TODO 请求的时候要设置 Accept-Encoding ： gzip, deflate, br
// TODO 请求的时候要设置 Referer : host url

type IPAddr struct {
	IP      string
	Address string // ip的地址信息
}

type SSLCertificateInfo struct {
	Url                   string `json:"url"`                   // url
	EffectiveTime         string `json:"effectiveTime"`         // 有效时间
	NotBefore             int64  `json:"notBefore"`             // 起始
	NotAfter              int64  `json:"notAfter"`              // 结束
	DNSName               string `json:"dnsName"`               // DNSName
	OCSPServer            string `json:"ocspServer"`            // OCSPServer
	CRLDistributionPoints string `json:"crlDistributionPoints"` // CRL分发点
	Issuer                string `json:"issuer"`                // 颁发者
	IssuingCertificateURL string `json:"issuingCertificateURL"` // 颁发证书URL
	PublicKeyAlgorithm    string `json:"publicKeyAlgorithm"`    // 公钥算法
	Subject               string `json:"subject"`               // 颁发对象
	Version               string `json:"version"`               // 版本
	SignatureAlgorithm    string `json:"signatureAlgorithm"`    // 证书算法
}

func (s SSLCertificateInfo) Echo() {
	txt := `
Url: %s 
有效时间: %s 
颁发对象: %s 
颁发者: %s 
颁发证书URL: %s 
公钥算法: %s 
证书算法: %s 
版本: %s 
DNSName: %s 
CRL分发点: %s 
OCSPServer: %s `
	log.Info(fmt.Sprintf(txt, s.Url, s.EffectiveTime, s.Subject, s.Issuer, s.IssuingCertificateURL, s.PublicKeyAlgorithm,
		s.SignatureAlgorithm, s.Version, s.DNSName, s.CRLDistributionPoints, s.OCSPServer))
}

func (s SSLCertificateInfo) Expire() int64 {
	return s.NotAfter - time.Now().Unix()
}

// DNS 信息
type DNSInfo struct {
	IPs           []string `json:"ips"`
	LookupCNAME   string   `json:"cname"`
	DnsServerIP   string   `json:"dnsServerIP"`
	DnsServerName string   `json:"dnsServerName"`
	IsCDN         bool     `json:"isCDN"`
	Ms            float64  `json:"ms"` // ms
}

// Whois 信息
type WhoisInfo struct {
	Root string
	Rse  string
}

type TDKI struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Icon        string `json:"icon"`
}
