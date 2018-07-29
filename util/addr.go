package util

import (
	"github.com/toolkits/net"
)

func GetLocalIP() string {
	ips, err := net.IntranetIP()
	if err != nil {
		return ""
	}

	// 虚拟机的本地网卡也会在这里列出，默认取第一个
	for _, ip := range ips {
		return ip
	}

	return ""
}
