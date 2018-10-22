package service

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

func GetDiscoveryAddr() string {
	return *flagDiscoveryAddr
}
