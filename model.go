package cellmesh

var (
	procName string
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
