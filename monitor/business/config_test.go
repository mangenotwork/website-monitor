package business

import "testing"

func TestAnalysisWebsiteSlowResponseTime(t *testing.T) {
	// master /data/website/:hostId 的实际返回结构
	raw := `{"code":0,"msg":"","data":{"host":"https://example.com","hostID":"h1","monitorRate":15,"websiteSlowResponseTime":4321}}`

	data := &Website{}
	if err := AnalysisData(raw, &data); err != nil {
		t.Fatal(err)
	}
	if data.WebsiteSlowResponseTime != 4321 {
		t.Fatalf("websiteSlowResponseTime=%d, want 4321", data.WebsiteSlowResponseTime)
	}

	// 兼容旧字段名（大写开头）
	legacy := &Website{}
	if err := AnalysisData(`{"code":0,"data":{"host":"https://example.com","WebsiteSlowResponseTime":3111}}`, &legacy); err != nil {
		t.Fatal(err)
	}
	if legacy.WebsiteSlowResponseTime != 3111 {
		t.Fatalf("legacy websiteSlowResponseTime=%d, want 3111", legacy.WebsiteSlowResponseTime)
	}
}
