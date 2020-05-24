package fx

var (
	// 节点名, 一般按进程名命名
	NodeName string

	// 服务分区, 一个分区下有多个分组
	NodeZone string

	// 服务分组, 同一台机器(IP), 分组相同
	NodeGroup string

	// 同类服务区分, 进程ID
	NodeIndex int

	// 本进程对应的SvcID
	LocalNodeID string

	// 公网IP
	WANIP string

	// 服务发现地址
	DiscoveryAddress string
)
