package dao

import (
	"io"
	"net"
	"regexp"
	"strings"

	"website-monitor/master/entity"

	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
)

var RootWhoisServers = "whois.iana.org:43"

func Whois(host string) *entity.WhoisInfo {
	hostList := strings.Split(host, ".")
	host = strings.Join(hostList[len(hostList)-2:len(hostList)], ".")
	info := &entity.WhoisInfo{}
	rootRse := whois(RootWhoisServers, host)
	info.Root = rootRse
	referList := regFindTxt(`(?is:refer:(.*?)\n)`, rootRse)
	if len(referList) > 0 {
		refer := utils.StrDeleteSpace(referList[0])
		rse := whois(refer+":43", host)
		info.Rse = rse
	}
	return info
}

func whois(server, host string) string {
	conn, _ := net.Dial("tcp", server)
	conn.Write([]byte(host + " \r\n"))
	buf := make([]byte, 1024*10)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Error(err)
	}
	rse := string(buf[:n])
	defer func() {
		conn.Close()
	}()
	return rse
}

func RegFindAll(regStr, rest string) [][]string {
	reg := regexp.MustCompile(regStr)
	List := reg.FindAllStringSubmatch(rest, -1)
	reg.FindStringSubmatch(rest)
	return List
}

// regFindTxt 执行正则提取 只取内容
func regFindTxt(regStr, txt string, property ...string) (dataList []string) {
	reg := regexp.MustCompile(regStr)
	resList := reg.FindAllStringSubmatch(txt, -1)
	for _, v := range resList {
		if len(v) < 1 {
			continue
		}
		if len(property) == 0 || strings.Count(v[0], strings.Join(property, " ")) > 0 {
			dataList = append(dataList, v[1])
		}
	}
	return
}
