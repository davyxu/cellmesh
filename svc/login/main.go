package main

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/link"
	_ "github.com/davyxu/cellmesh/svc/login/verify"
)

func main() {
	cellmesh.Init("login")
	cellmesh.LogParameter()
	cellmesh.ConnectDiscovery()

	// 网关对客户端连接
	link.StartService(&link.ServiceParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "tcp.client",
		SvcName:       "verify",
		ListenAddress: ":8001",
		Queue:         cellmesh.Queue,
		EventCallback: fx.MakeIOCEventHandler(fx.MessageRegistry),
	})

	// 服务互联
	link.LinkService(&link.ServiceParameter{
		PeerType: "tcp.Connector",
		NetProc:  "tcp.svc",
		SvcName:  "hub",
		Queue:    cellmesh.Queue,
	})

	link.CheckReady()

	cellmesh.WaitExitSignal()

}
