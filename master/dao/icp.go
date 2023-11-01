package dao

import (
	"fmt"
	"website-monitor/master/entity"
)

var GetICPUrl = func(host string) string { return fmt.Sprintf("https://www.beianx.cn/search/%s", host) }

// GetICP 获取备案信息  TODO BUG...
func GetICP(host string) *entity.ICPInfo {
	ipcInfo := &entity.ICPInfo{
		Host: host,
	}
	//// 有反爬，第一次返回的html是js生成acw_sc__v2 set cookie, 然后刷新
	//h := gt.Header{"Cookie": "acw_sc__v2=654203242be7b8650bff3923d72a260cc6055a87;"}
	//ctx, _ := gt.Get(GetICPUrl(host), h)
	//rseList, _ := gt.GetPointIDHTML(ctx.Html, "div", "aa_result")
	//rse := gt.RegHtmlTbody(rseList[0])
	//log.Error("rse = ", rse)
	//if len(rse) > 0 {
	//	tdList := gt.RegHtmlTdTxt(rse[0])
	//	if len(tdList) < 8 {
	//		return ipcInfo
	//	}
	//	td1 := tdList[1]
	//	td1Rse := gt.RegHtmlATxt(td1)
	//	if len(td1Rse) > 0 {
	//		ipcInfo.Company = td1Rse[0]
	//	}
	//	td2 := tdList[2]
	//	ipcInfo.Nature = utils.CleaningStr(td2)
	//	td3 := tdList[3]
	//	ipcInfo.IPC = utils.CleaningStr(td3)
	//	td4 := tdList[4]
	//	ipcInfo.WebsiteName = utils.CleaningStr(td4)
	//	td5 := tdList[5]
	//	td5Rse := gt.RegHtmlATxt(td5)
	//	if len(td5Rse) > 0 {
	//		ipcInfo.WebsiteIndex = td5Rse[0]
	//	}
	//	td6 := tdList[6]
	//	td6Rse := gt.RegHtmlDivTxt(td6)
	//	if len(td6Rse) > 0 {
	//		ipcInfo.AuditDate = td6Rse[0]
	//	}
	//	td7 := tdList[7]
	//	ipcInfo.RestrictAccess = utils.CleaningStr(td7)
	//}
	return ipcInfo
}
