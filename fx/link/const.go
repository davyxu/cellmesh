package link

import (
	"github.com/davyxu/cellnet"
)

type NodeParameter struct {
	NodeName string // 服务名,注册到服务发现，本进程服务可以填空
	PeerType string // cellnet的PeerType
	NetProc  string // cellnet处理器名称

	ListenAddress string // socket侦听地址, 仅限开启服务使用
	Queue         cellnet.EventQueue
	EventCallback cellnet.EventCallback
}
