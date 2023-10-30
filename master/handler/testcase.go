package handler

import (
	"github.com/mangenotwork/common/ginHelper"
	"website-monitor/master/constname"
	"website-monitor/master/dao"
)

func NoticeUpdateWebsiteTest(c *ginHelper.GinCtx) {
	dao.Servers.NoticeAll(constname.NoticeUpdateWebsiteLabel, []byte("test"), dao.Servers.SetNoticeRetry(2, 3000))
}

func NoticeUpdateWebsiteAllUrlTest(c *ginHelper.GinCtx) {
	dao.Servers.NoticeAll(constname.NoticeUpdateWebsiteAllUrlLabel, []byte("test"), dao.Servers.SetNoticeRetry(2, 3000))
}

func NoticeUpdateWebsitePointTest(c *ginHelper.GinCtx) {
	dao.Servers.NoticeAll(constname.NoticeUpdateWebsitePointLabel, []byte("test"), dao.Servers.SetNoticeRetry(2, 3000))
}
