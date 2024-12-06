package dao

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mangenotwork/common/log"
	gt "github.com/mangenotwork/gathertool"
)

type IPInfo struct {
	IP      string
	Address string
}

func SetIP(ipInfo *IPInfo) error {
	return DB.Set(IPTable, ipInfo.IP, ipInfo.Address)
}

func GetIP(ip string) string {
	ipList := strings.Split(ip, ":")
	if len(ipList) > 0 {
		ip = ipList[0]
	}

	address := ""
	err := DB.Get(IPTable, ip, &address)
	if err != nil || address == "" {
		req := ReqIP(ip)
		address = req.Address
	}

	return address
}

const GetMyIPInfoUrl = "https://www.ip.cn/api/index?ip=&type=0"
const GetIPInfoUrl = "https://www.ip.cn/api/index?ip=%s&type=1"

var MyIP *IPInfo
var MyIPOnce sync.Once

func GetMyIP() *IPInfo {
	MyIPOnce.Do(func() {
		MyIP = GetNativeIP()
	})
	return MyIP
}

func GetNativeIP() *IPInfo {
	ctx, _ := gt.Get(GetMyIPInfoUrl)
	ip, _ := gt.JsonFind2Str(ctx.Json, "/ip")
	address, _ := gt.JsonFind2Str(ctx.Json, "/address")
	ipInfo := &IPInfo{
		IP:      ip,
		Address: address,
	}

	_ = SetIP(ipInfo)
	return ipInfo
}

func ReqIP(ipStr string) *IPInfo {
	ctx, _ := gt.Get(fmt.Sprintf(GetIPInfoUrl, ipStr))
	ip, _ := gt.JsonFind2Str(ctx.Json, "/ip")
	log.Info("ip = ", ip)

	address, _ := gt.JsonFind2Str(ctx.Json, "/address")
	log.Info("address = ", address)
	ipInfo := &IPInfo{
		IP:      ip,
		Address: address,
	}

	_ = SetIP(ipInfo)
	return ipInfo
}
