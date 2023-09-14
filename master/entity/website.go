package entity

import (
	"fmt"
	"github.com/mangenotwork/common/log"
	"time"
)

// Website 监测网站
type Website struct {
	Host string // 默认取根url进行检测，主键

	ID int64 // 网站ID, 独立id

	// 监测网络频率  单位 ms
	MonitorRate int64

	// 扫描站点深度 默认 2
	UriDepth int64

	// 扫描网站频率  单位秒
	ScanRate int64

	// 扫描网站内容
	ScanCheckUp ScanCheckUp

	// 设置报警规则
	AlarmRule AlarmRule

	Created int64

	Info WebsiteInfo
}

// WebsiteInfo 网站基本信息
type WebsiteInfo struct {
	Title              string
	Description        string
	Keywords           string
	Icon               string
	DNS                string             // 网站DNS信息
	IPAddr             []*IPAddr          // 网站IP和属地
	Server             string             // response Headers 读取 headers Server
	ContentEncoding    string             //  response Headers 是否压缩 读取 headers Content-Encoding
	ContentLanguage    string             // response Headers 语言 读取 headers Content-Language
	CDN                string             // 网站cdn信息
	SSLCertificateInfo SSLCertificateInfo // 证书信息
	Filing             string             // 网站备案信息
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

// AlarmRule 报警规则
type AlarmRule struct {
	WebsiteSlowResponseTime  int64 // 单位ms  网站响应有多慢才记录报警
	WebsiteSlowResponseCount int64 // 连续几次慢就发送邮件通知
	SSLCertificateExpire     int64 // 单位天 证书还有几天过期就触发报警
	NotTDK                   bool  // 在随机监测中发现网站页面没有存在 title, description, keywords 就触发报警
	BadLink                  bool  // 扫描的时候存在死链就报警
	ExtLinkChange            bool  // 扫描间隔判断外链有变化则报警,可以监测到是否被劫持
}

// ScanCheckUp 网站扫描检查内容
type ScanCheckUp struct {
	ScanExtLinks bool // 是否检查外链,对比上一次扫描数据判别外链是否变化
	ScanBadLink  bool // 扫描死链接
	// TODO... 安全扫描， Sql注入, XSS等等...
}

// WebSiteUrl 网站的URL存储
type WebSiteUrl struct {
	HostID  string
	AllUri  []string // 抓取到的所有链接
	ExtLink []string // 外链
	BadLink []string // 死链
	JsLink  []string // 资源文件 js
	CssLink []string // 资源文件 css
}
