package business

import "sync"

var WebsitePointDataMap sync.Map

// WebSitePoint 指定监测点(URL)
type WebSitePoint struct {
	HostID string   `json:"hostID"`
	URL    []string `json:"url"`
}

func GetWebSitePoint(hostId string) ([]string, int64) {
	value, ok := WebsitePointDataMap.Load(hostId)
	if !ok {
		return []string{}, 0
	}
	data := value.(*WebSitePoint)
	return data.URL, int64(len(data.URL))
}
