package business

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestRequestNormal(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	code, ms, err := request(srv.URL)
	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	if code != http.StatusOK {
		t.Fatalf("code=%d", code)
	}
	if ms < 0 {
		t.Fatalf("ms=%d", ms)
	}
}

// 服务端声明 Content-Length 100，只写 5 字节就断开。
// gathertool 会忽略 io.ErrUnexpectedEOF，这里必须返回错误。
func TestRequestUnexpectedEOF(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "no hijack", http.StatusInternalServerError)
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			return
		}
		defer conn.Close()
		_, _ = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 100\r\n\r\nshort"))
	}))
	defer srv.Close()

	code, _, err := request(srv.URL)
	if err == nil {
		t.Fatalf("截断响应未返回错误, code=%d", code)
	}
	if !strings.Contains(err.Error(), "unexpected EOF") {
		t.Fatalf("err=%v, 期望包含 unexpected EOF", err)
	}
}

func TestIsTimeoutError(t *testing.T) {
	if !isTimeoutError(context.DeadlineExceeded) {
		t.Fatal("context.DeadlineExceeded 应判为超时")
	}
	if !isTimeoutError(&url.Error{Op: "Get", URL: "http://x", Err: context.DeadlineExceeded}) {
		t.Fatal("url.Error 包装的 deadline 应判为超时")
	}
	if !isTimeoutError(&net.DNSError{IsTimeout: true}) {
		t.Fatal("net.DNSError timeout 应判为超时")
	}
	if isTimeoutError(io.ErrUnexpectedEOF) {
		t.Fatal("unexpected EOF 不应判为超时")
	}
	if isTimeoutError(nil) {
		t.Fatal("nil 不应判为超时")
	}
}

func TestRequestTimeout(t *testing.T) {
	old := monitorHTTPClient
	monitorHTTPClient = &http.Client{Timeout: 100 * time.Millisecond}
	defer func() { monitorHTTPClient = old }()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
	}))
	defer srv.Close()

	_, _, err := request(srv.URL)
	if err == nil {
		t.Fatal("超时未返回错误")
	}
	if !isTimeoutError(err) {
		t.Fatalf("err=%v, 应判为超时", err)
	}
}
