package dao

import (
	"fmt"
	"io"
	"net"
	"regexp"
)

var RootWhoisServers = "whois.iana.org:43"

func Whois(host string) {
	conn, _ := net.Dial("tcp", RootWhoisServers)
	conn.Write([]byte("github.com \r\n"))
	buf := make([]byte, 1024*10)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		fmt.Println(err)
	}
	rse := string(buf[:n])
	fmt.Println(rse)
	RegFindAll(`(?is:refer:.*?\\n)`, rse)
	conn.Close()
}

func RegFindAll(regStr, rest string) [][]string {
	reg := regexp.MustCompile(regStr)
	List := reg.FindAllStringSubmatch(rest, -1)
	reg.FindStringSubmatch(rest)
	return List
}
