package service

import (
	"fmt"
)

var (
	procName  string
	LinkRules []MatchRule // 互联发现规则
)

// 获取当前服务进程名称
func GetProcName() string {
	return procName
}

// 获取外网IP
func GetWANIP() string {
	return *flagWANIP
}

func GetSvcGroup() string {
	return *flagSvcGroup
}

func GetSvcIndex() int {
	return *flagSvcIndex
}

// 构造服务ID
func MakeSvcID(svcName string, svcIndex int, svcGroup string) string {
	return fmt.Sprintf("%s#%d@%s", svcName, svcIndex, svcGroup)
}

// 构造指定服务的ID
func MakeLocalSvcID(svcName string) string {
	return MakeSvcID(svcName, *flagSvcIndex, *flagSvcGroup)
}

func GetLocalSvcID() string {
	return MakeLocalSvcID(GetProcName())
}

func GetDiscoveryAddr() string {
	return *flagDiscoveryAddr
}
