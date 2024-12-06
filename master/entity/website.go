package entity

import (
	"fmt"
	"time"

	"github.com/mangenotwork/common/log"
)

// Website 监测网站
type Website struct {
	// 默认取根url进行检测，主键
	Host string `json:"host"`

	HostID string `json:"hostID"`

	// 监测频率  单位 ms
	MonitorRate int64 `json:"monitorRate"`

	// 保证监测的有效性:每次监测前用过对照组+ping来确认当前监测器网络的情况，只有当网络情况好才执行监测请求
	// 对照组Url
	ContrastUrl string `json:"contrastUrl"`

	// 对照组响应超过这个时间判定为当前网络不稳定不执行本次监测请求
	ContrastTime int64 `json:"contrastTime"`

	// 用于检查当前网络
	Ping string `json:"ping"`

	// ping响应超过这个时间判定为当前网络不稳定不执行本次监测请求
	PingTime int64 `json:"pingTime"`

	Notes   string `json:"notes"`
	Created int64  `json:"created"`
}

// WebsiteAlarmRule 报警规则
type WebsiteAlarmRule struct {
	Host                     string `json:"host"`                     // 主键与website Host对应
	HostID                   string `json:"hostID"`                   // key 键
	WebsiteSlowResponseTime  int64  `json:"websiteSlowResponseTime"`  // 单位ms  网站响应有多慢才记录报警
	WebsiteSlowResponseCount int64  `json:"websiteSlowResponseCount"` // 连续几次慢就发送邮件通知
	SSLCertificateExpire     int64  `json:"SSLCertificateExpire"`     // 单位天 证书还有几天过期就触发报警
	NotTDK                   bool   `json:"notTDK"`                   //  在随机监测中发现网站页面没有存在 title, description, keywords 就触发报警
	BadLink                  bool   `json:"badLink"`                  // 扫描的时候存在死链就报警
	ExtLinkChange            bool   `json:"extLinkChange"`            // 扫描间隔判断外链有变化则报警,可以监测到是否被劫持
}

// WebsiteScanCheckUp 网站扫描检查内容
type WebsiteScanCheckUp struct {
	Host         string `json:"host"`         // 主键与website Host对应
	HostID       string `json:"hostID"`       // key 键
	ScanDepth    int64  `json:"uriDepth"`     // 扫描站点深度 默认 2
	ScanRate     int64  `json:"scanRate"`     // 扫描网站频率  单位秒
	ScanExtLinks bool   `json:"scanExtLinks"` // 是否检查外链,对比上一次扫描数据判别外链是否变化
	ScanBadLink  bool   `json:"scanBadLink"`  // 是否扫描死链接
	ScanTDK      bool   `json:"scanTDK"`      // 是否扫描进行TDK检查
	// TODO... 安全扫描， Sql注入, XSS等等...
}

// WebsiteInfo 网站基本信息
type WebsiteInfo struct {
	Host               string              `json:"host"`               // 主键与website Host对应
	HostID             string              `json:"hostID"`             // key 键
	Title              string              `json:"title"`              // Title
	Description        string              `json:"description"`        // Description
	Keywords           string              `json:"keywords"`           // Keywords
	Icon               string              `json:"icon"`               // Icon
	DNS                *DNSInfo            `json:"DNS"`                // 网站DNS信息
	IPAddr             []*IPAddr           `json:"IPAddr"`             // 网站IP和属地
	Server             string              `json:"server"`             // response Headers 读取 headers Server
	ContentEncoding    string              `json:"contentEncoding"`    //  response Headers 是否压缩 读取 headers Content-Encoding
	ContentLanguage    string              `json:"contentLanguage"`    // response Headers 语言 读取 headers Content-Language
	SSLCertificateInfo *SSLCertificateInfo `json:"SSLCertificateInfo"` // 证书信息
	Filing             string              `json:"filing"`             // 网站备案信息
	Whois              *WhoisInfo          `json:"whois"`              // Whois信息
	IPC                *ICPInfo            `json:"IPC"`                // ipc 信息
}

// WebSiteUrl 网站的URL存储
type WebSiteUrl struct {
	Host    string            `json:"host"`    // 主键与website Host对应
	HostID  string            `json:"hostID"`  // key 键
	AllUri  []string          `json:"allUri"`  // 抓取到的所有链接
	ExtLink []string          `json:"extLink"` // 外链
	BadLink []string          `json:"badLink"` // 死链
	NoneTDK map[string]string `json:"noneTDK"` // 空TDK的链接
	JsLink  []string          `json:"jsLink"`  // 资源文件 js
	CssLink []string          `json:"cssLink"` // 资源文件 css
	ImgLink []string          `json:"imgLink"` // 采集到的图片文件链接
}

// WebSiteUrlPoint 指定网站监测Url
type WebSiteUrlPoint struct {
	Host   string   `json:"host"`   // 主键与website Host对应
	HostID string   `json:"hostID"` // key 键
	Url    []string `json:"url"`
}

// WebSiteAlert 监控报警信息
type WebSiteAlert struct {
	Host   string       `json:"host"`   // 主键与website Host对应
	HostID string       `json:"hostID"` // key 键
	List   []*AlertData `json:"list"`
}

type IPAddr struct {
	IP      string `json:"ip"`
	Address string `json:"address"` // ip的地址信息
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

// DNSInfo 信息
type DNSInfo struct {
	IPs           []string `json:"ips"`
	LookupCNAME   string   `json:"cname"`
	DnsServerIP   string   `json:"dnsServerIP"`
	DnsServerName string   `json:"dnsServerName"`
	IsCDN         bool     `json:"isCDN"`
	Ms            float64  `json:"ms"` // ms
}

// WhoisInfo Whois 信息
type WhoisInfo struct {
	Root string `json:"root"`
	Rse  string `json:"rse"`
}

type TDKI struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Icon        string `json:"icon"`
}

// ICPInfo icp信息
type ICPInfo struct {
	Host           string `json:"host"`           // 网站
	Company        string `json:"company"`        // 公司
	Nature         string `json:"nature"`         // 性质
	IPC            string `json:"ipc"`            // ipc
	WebsiteName    string `json:"websiteName"`    // 网站名称
	WebsiteIndex   string `json:"websiteIndex"`   // 网站主页
	AuditDate      string `json:"auditDate"`      // 审核日期
	RestrictAccess string `json:"restrictAccess"` // 是否限制接入
}

// WebSitePoint 指定监测点(URL)
type WebSitePoint struct {
	HostID string   `json:"hostID"`
	URL    []string `json:"url"`
}
