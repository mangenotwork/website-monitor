package handler

import (
	"fmt"
	"net"
	"os"
	"runtime"
)

type OSInfo struct {
	HostName      string
	OSType        string
	OSArch        string
	CpuCoreNumber string
	InterfaceInfo string
}

// GetHostName 获取host 命名
func GetHostName() string {
	name, err := os.Hostname()
	if err != nil {
		name = "null"
	}
	return name
}

// GetSysType 获取host 系统类型
func GetSysType() string {
	return runtime.GOOS
}

// GetSysArch 获取系统架构
func GetSysArch() string {
	return runtime.GOARCH
}

// GetCpuCoreNumber 获取cpu核心数
func GetCpuCoreNumber() string {
	return fmt.Sprintf("%d核", runtime.GOMAXPROCS(0))
}

func GetInterfaceInfo() string {
	rse := ""
	iFaces, _ := net.Interfaces()

	for _, v := range iFaces {
		addr, err := v.Addrs()
		if err != nil {
			continue
		}

		for _, a := range addr {
			ipNet, ok := a.(*net.IPNet)

			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				rse += fmt.Sprintf("网卡:%s;ip:%s;mac:%s;\n", v.Name, ipNet.IP, v.HardwareAddr.String())
			}

		}

	}

	return rse
}
