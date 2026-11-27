package dao

import (
	"net"
	"net/url"
	"strings"
)

// NormalizeDomain 从 URL / 域名里取出不含协议、端口、路径的 host。
// 例：https://www.doubao.com/ -> www.doubao.com
func NormalizeDomain(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if !strings.Contains(s, "://") {
		s = "http://" + s
	}

	if u, err := url.Parse(s); err == nil && u.Hostname() != "" {
		return strings.ToLower(strings.TrimSuffix(u.Hostname(), "."))
	}

	// 非法 URL 的退化处理
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")
	if i := strings.IndexAny(s, "/?#"); i >= 0 {
		s = s[:i]
	}
	if h, _, err := net.SplitHostPort(s); err == nil {
		s = h
	}
	return strings.ToLower(strings.Trim(s, ". "))
}

// 常见二级后缀，whois 需要的是可注册域名（eTLD+1）
var multiPartSuffixes = map[string]struct{}{
	"com.cn": {}, "net.cn": {}, "org.cn": {}, "gov.cn": {}, "edu.cn": {}, "ac.cn": {},
	"com.hk": {}, "com.mo": {}, "com.tw": {}, "co.uk": {}, "org.uk": {}, "co.jp": {},
	"com.au": {}, "co.nz": {},
}

// rootDomain 取可注册域名：www.example.com.cn -> example.com.cn
func rootDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) <= 2 {
		return domain
	}

	suffix := strings.Join(parts[len(parts)-2:], ".")
	if _, ok := multiPartSuffixes[suffix]; ok && len(parts) >= 3 {
		return strings.Join(parts[len(parts)-3:], ".")
	}
	return suffix
}
