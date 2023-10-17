package business

import (
	"github.com/mangenotwork/common/log"
	gt "github.com/mangenotwork/gathertool"
	"sync"
)

// WebsiteUrlDataMap 监测网站扫描到的url
var WebsiteUrlDataMap sync.Map

func SetWebsiteUrlData(hostId string) {
	ctx, err := gt.Get(MasterHTTP + GetWebsiteAllUrlAPI(hostId))
	if err != nil {
		log.Error(err)
		return
	}
	allUrl := make([]string, 0)
	err = AnalysisData(ctx.Json, &allUrl)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("allUrl = ", allUrl)
	WebsiteUrlDataMap.Store(hostId, allUrl)
}

func GetWebsiteUrlDataMap(hostId string) []string {
	v, ok := WebsiteUrlDataMap.Load(hostId)
	if !ok {
		return make([]string, 0)
	}
	return v.([]string)
}
