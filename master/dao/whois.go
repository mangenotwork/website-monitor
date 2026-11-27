package dao

import (
	"io"
	"net"
	"regexp"
	"strings"
	"time"

	"website-monitor/master/entity"

	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
)

var RootWhoisServers = "whois.iana.org:43"

func Whois(host string) *entity.WhoisInfo {
	info := &entity.WhoisInfo{}

	host = rootDomain(NormalizeDomain(host))
	if host == "" {
		return info
	}

	rootRse := whois(RootWhoisServers, host)
	info.Root = rootRse
	referList := regFindTxt(`(?is:refer:(.*?)\n)`, rootRse)

	if len(referList) > 0 {
		refer := strings.TrimSpace(utils.StrDeleteSpace(referList[0]))
		if refer != "" {
			info.Rse = whois(refer+":43", host)
		}
	}

	return info
}

func whois(server, host string) string {
	conn, err := net.DialTimeout("tcp", server, 8*time.Second)
	if err != nil {
		log.Error("whois 连接失败 ", server, " err: ", err)
		return ""
	}
	defer func() {
		_ = conn.Close()
	}()

	_ = conn.SetDeadline(time.Now().Add(15 * time.Second))

	if _, err := conn.Write([]byte(host + "\r\n")); err != nil {
		log.Error("whois 写入失败 ", server, " err: ", err)
		return ""
	}

	buf, err := io.ReadAll(io.LimitReader(conn, 1<<20))
	if err != nil && err != io.EOF {
		log.Error("whois 读取失败 ", server, " err: ", err)
	}

	return string(buf)
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
