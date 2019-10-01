package cellmesh

import (
	"fmt"
	"github.com/davyxu/cellnet/util"
	"net"
	"os"
	"strings"
)

var (
	thisSvcID string
)

// ip+PID的16进制数值字符串，每次启动变化
func netProcID() string {

	// 一次启动不会变化
	if thisSvcID != "" {
		return thisSvcID
	}

	// 兼容ipv6
	ipParts := net.ParseIP(util.GetLocalIP())

	var sb strings.Builder
	for _, p := range ipParts {
		if p == 0 || p == 255 {
			continue
		}

		sb.WriteString(fmt.Sprintf("%x", p))
	}

	sb.WriteString(fmt.Sprintf("%x", os.Getpid()))

	thisSvcID = sb.String()

	return thisSvcID
}

// 全局唯一的svcid 格式:  svcName@netProcID

// 构造指定服务的ID
func MakeSvcID(svcName string) string {
	return svcName + "@" + netProcID()
}

// 获得本进程的服务id
func GetLocalSvcID() string {
	return MakeSvcID(GetProcName())
}
