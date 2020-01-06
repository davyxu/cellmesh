package link

import (
	"github.com/davyxu/cellnet"
)

type ServiceParameter struct {
	SvcName  string // 服务名,注册到服务发现，本进程服务可以填空
	PeerType string // cellnet的PeerType
	NetProc  string // cellnet处理器名称

	ListenAddress string // socket侦听地址, 仅限开启服务使用
	Queue         cellnet.EventQueue
	EventCallback cellnet.EventCallback
}

func (self *ServiceParameter) MakeServiceDefault() {

	if self.PeerType == "" {
		self.PeerType = "tcp.Acceptor"
	}

	if self.NetProc == "" {
		self.NetProc = "tcp.svc"
	}
}

func (self *ServiceParameter) MakeConnectorDefault() {
	if self.PeerType == "" {
		self.PeerType = "tcp.Connector"
	}
	if self.NetProc == "" {
		self.NetProc = "tcp.svc"
	}
}
