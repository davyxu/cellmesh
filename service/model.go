package service

var (
	matchNodes []string
	procName   string
)

// 获取当前服务进程名称
func GetProcName() string {
	return procName
}

// 获取当前节点
func GetNode() string {
	return *flagNode
}

// 获取外网IP
func GetWANIP() string {
	return *flagWANIP
}

// 获取要匹配节点名
func GetMatchNodes() []string {
	return matchNodes
}

// 构造服务ID
func MakeServiceID(svcName string) string {
	return svcName + "@" + GetNode()
}
