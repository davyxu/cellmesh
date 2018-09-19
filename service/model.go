package service

var (
	procName string
)

// 获取当前服务进程名称
func GetProcName() string {
	return procName
}

// 获取当前节点(服务侦听用)
func GetNode() string {
	return *flagNode
}

// 获取外网IP
func GetWANIP() string {
	return *flagWANIP
}

// 构造服务ID
func MakeServiceID(svcName string) string {
	return svcName + "@" + GetNode()
}
