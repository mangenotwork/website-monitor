package dao

import (
	"github.com/mangenotwork/common/log"
	gt "github.com/mangenotwork/gathertool"
	"math/rand"
	"strings"
	"time"
	"website-monitor/master/entity"
)

type HostScanUrl struct {
	Host     string
	Depth    int64               // 页面深度
	UrlSet   map[string]struct{} // 采集到host下的链接
	ExtLinks map[string]struct{} // 采集到外链
	BadLinks map[string]struct{} // 采集到死链接
	NoneTDK  map[string]string   // 检查空tdk
	Count    int64
	MaxCount int64
}

func Scan(host, id string, depth int64) {
	hostScan := &HostScanUrl{
		Host:     host,
		Depth:    depth,
		UrlSet:   make(map[string]struct{}),
		ExtLinks: make(map[string]struct{}),
		BadLinks: make(map[string]struct{}),
		NoneTDK:  make(map[string]string),
	}
	hostScan.Run()
	websiteUrl := &entity.WebSiteUrl{
		Host:    host,
		HostID:  id,
		AllUri:  make([]string, 0),
		ExtLink: make([]string, 0),
		BadLink: make([]string, 0),
		NoneTDK: make(map[string]string),
		JsLink:  make([]string, 0),
		CssLink: make([]string, 0),
	}
	log.Info("扫描完成....")
	for u, _ := range hostScan.UrlSet {
		websiteUrl.AllUri = append(websiteUrl.AllUri, u)
	}
	for e, _ := range hostScan.ExtLinks {
		websiteUrl.ExtLink = append(websiteUrl.ExtLink, e)
	}
	for b, _ := range hostScan.BadLinks {
		websiteUrl.BadLink = append(websiteUrl.BadLink, b)
	}
	for k, v := range hostScan.NoneTDK {
		websiteUrl.NoneTDK[k] = v
	}

	err := DB.Set(WebSiteURITable, id, websiteUrl)
	if err != nil {
		log.Error("保存扫描的数据失败 err = ", err)
	}
}

func (scan *HostScanUrl) Run() {
	scan.do(scan.Host, 0)
}

func (scan *HostScanUrl) do(caseUrl string, df int64) {
	if len(caseUrl) < 1 {
		return
	}
	if df > scan.Depth {
		return
	}
	// 如果不是host下的域名,就是外链
	if strings.Index(caseUrl, scan.Host) == -1 {
		if string(caseUrl[0]) == "/" {
			caseUrl = scan.Host + caseUrl
			goto G
		} else if string(caseUrl[0]) != "/" && string(caseUrl[0]) != "#" {
			scan.ExtLinks[caseUrl] = struct{}{}
		}
		return
	}
G:
	if _, ok := scan.UrlSet[caseUrl]; ok {
		return
	}
	log.Info("执行扫描网站 --> ", caseUrl)
	sleepMs(500, 2500)
	ctx, err := gt.Get(caseUrl)
	if err != nil {
		log.Error(err)
		return
	}
	// 记录死链接
	if ctx.StateCode >= 400 {
		scan.BadLinks[caseUrl] = struct{}{}
	}
	// 检查空TDK
	contentType := ctx.Resp.Header.Get("Content-Type")
	if strings.Index(contentType, "text/html") != -1 {
		rse := checkTDK(ctx.Html)
		if len(rse) > 0 {
			scan.NoneTDK[caseUrl] = rse
		}
	}
	df++
	scan.UrlSet[caseUrl] = struct{}{}
	scan.Count++
	a := gt.RegHtmlA(ctx.Html)
	for _, v := range a {
		links := gt.RegHtmlHrefTxt(v)
		if len(links) < 1 {
			continue
		}
		link := links[0]
		// 请求并验证
		scan.do(link, df)
	}
}

func checkTDK(html string) string {
	rse := ""
	title := gt.RegHtmlTitleTxt(html)
	if len(title) < 1 {
		rse += "title:none;"
	}
	description := gt.RegHtmlDescriptionTxt(html)
	if len(description) < 1 {
		rse += "description:none;"
	}
	keyword := gt.RegHtmlKeywordTxt(html)
	if len(keyword) < 1 {
		rse += "keyword:none;"
	}
	return rse
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var randEr = rand.New(rand.NewSource(time.Now().UnixNano()))

// SetSleepMs 设置请求随机休眠时间， 单位毫秒
func sleepMs(min, max int) {
	r := randEr.Intn(max) + min
	time.Sleep(time.Duration(r) * time.Millisecond)
}