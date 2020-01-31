package cellmesh

import "github.com/davyxu/cellnet"

var (
	// 全局队列
	Queue cellnet.EventQueue

	// 进程名
	ProcName string

	// 公网IP
	WANIP string

	// 服务发现地址
	DiscoveryAddress string
)
