package dao

import (
	"fmt"
	"github.com/mangenotwork/common/utils"
	gt "github.com/mangenotwork/gathertool"
	"website-monitor/master/entity"
)

var GetICPUrl = func(host string) string { return fmt.Sprintf("https://www.beianx.cn/search/%s", host) }

// GetICP 获取备案信息
func GetICP(host string) *entity.ICPInfo {
	ipcInfo := &entity.ICPInfo{
		Host: host,
	}
	ctx, _ := gt.Get(GetICPUrl(host))
	rse := gt.RegHtmlTbody(ctx.Html)
	if len(rse) > 0 {
		tdList := gt.RegHtmlTdTxt(rse[0])
		if len(tdList) < 8 {
			return ipcInfo
		}
		td1 := tdList[1]
		td1Rse := gt.RegHtmlATxt(td1)
		if len(td1Rse) > 0 {
			ipcInfo.Company = td1Rse[0]
		}
		td2 := tdList[2]
		ipcInfo.Nature = utils.CleaningStr(td2)
		td3 := tdList[3]
		ipcInfo.IPC = utils.CleaningStr(td3)
		td4 := tdList[4]
		ipcInfo.WebsiteName = utils.CleaningStr(td4)
		td5 := tdList[5]
		td5Rse := gt.RegHtmlATxt(td5)
		if len(td5Rse) > 0 {
			ipcInfo.WebsiteIndex = td5Rse[0]
		}
		td6 := tdList[6]
		td6Rse := gt.RegHtmlDivTxt(td6)
		if len(td6Rse) > 0 {
			ipcInfo.AuditDate = td6Rse[0]
		}
		td7 := tdList[7]
		ipcInfo.RestrictAccess = utils.CleaningStr(td7)
	}
	return ipcInfo
}
