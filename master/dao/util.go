package dao

func urlStr(url string) string {
	l := len(url)
	if l < 1 {
		panic("url is null")
	}
	if l > 8 && (url[:7] == "http://" || url[:8] == "https://") {
		return url
	}
	return "https://" + url
}

// InspectCode 请求状态码出现 >=400, >=500 的状态码视为网站url有问题
func InspectCode(code int) bool {
	if code >= 400 || code >= 500 {
		return false
	}
	return true
}
