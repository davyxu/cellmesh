package main

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/link"
	_ "github.com/davyxu/cellmesh/svc/login/verify"
)

func main() {
	fx.Init("login")
	fx.LogParameter()
	fx.ConnectDiscovery()

	// 网关对客户端连接
	link.ListenService(&link.ServiceParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "tcp.client",
		SvcName:       "verify",
		ListenAddress: ":8001",
		Queue:         fx.Queue,
		EventCallback: fx.MakeIOCEventHandler(fx.MessageRegistry),
	})

	// 服务互联
	link.ConnectService(&link.ServiceParameter{
		PeerType: "tcp.Connector",
		NetProc:  "tcp.svc",
		SvcName:  "hub",
		Queue:    fx.Queue,
	})

	link.CheckReady()

	fx.WaitExitSignal()

}
