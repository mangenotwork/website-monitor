package dao

import "testing"

func TestNormalizeDomain(t *testing.T) {
	cases := map[string]string{
		"https://www.doubao.com/":        "www.doubao.com",
		"http://www.doubao.com":          "www.doubao.com",
		"www.doubao.com:443/path?a=1":    "www.doubao.com",
		"HTTP://Example.COM./":           "example.com",
		"https://www.example.com.cn/a/b": "www.example.com.cn",
		"":                               "",
		"https://":                       "",
	}
	for in, want := range cases {
		if got := NormalizeDomain(in); got != want {
			t.Errorf("NormalizeDomain(%q)=%q, want %q", in, got, want)
		}
	}
}

func TestRootDomain(t *testing.T) {
	cases := map[string]string{
		"www.doubao.com":     "doubao.com",
		"www.example.com.cn": "example.com.cn",
		"a.b.example.co.uk":  "example.co.uk",
		"localhost":          "localhost",
		"example.com":        "example.com",
	}
	for in, want := range cases {
		if got := rootDomain(in); got != want {
			t.Errorf("rootDomain(%q)=%q, want %q", in, got, want)
		}
	}
}

// 复现旧 panic：whois 连不上时 conn 为 nil，不能再 Write
func TestWhoisDialErrorNoPanic(t *testing.T) {
	old := RootWhoisServers
	RootWhoisServers = "127.0.0.1:1" // 必然拒绝连接
	defer func() { RootWhoisServers = old }()

	info := Whois("https://www.doubao.com/")
	if info == nil {
		t.Fatal("Whois 返回 nil")
	}
	if info.Root != "" || info.Rse != "" {
		t.Fatalf("连接失败时应返回空结果: %+v", info)
	}
}
