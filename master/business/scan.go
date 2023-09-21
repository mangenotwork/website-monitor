package business

import (
	"time"
)

func Collect() {
	timer := time.NewTimer(time.Hour * 1) //初始化定时器
	for {
		select {
		case <-timer.C:
			//log.Info("执行采集...")
			// 获取站点对象
			//WebSiteObj.Range(func(key any, value any) bool {
			//	web := value.(*WebSiteItem)
			//	webSiteUri := model.NewWebSiteUri(web.ID)
			//	webSiteUri.Collect(web.HealthUri, int(web.UriDepth))
			//	return true
			//})
			timer.Reset(time.Hour * 1)
		}
	}
}
