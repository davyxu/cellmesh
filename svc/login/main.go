package main

import (
	_ "github.com/davyxu/cellmesh/fx/proc"
	"github.com/davyxu/cellmesh/link"
	"time"
)

import (
	"github.com/davyxu/cellmesh/fx"
	_ "github.com/davyxu/cellmesh/svc/login/verify"
)

func main() {
	fx.Init("login")
	fx.LogParameter()
	link.ConnectDiscovery()

	link.MonitorService("frontend", time.Second*3)
	link.MonitorService("game", time.Second*3)

	// 网关对客户端连接
	link.ListenNode(&link.NodeParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "tcp.client",
		SvcName:       "verify",
		ListenAddress: ":8001",
		EventCallback: fx.MakeIOCEventHandler(fx.MessageRegistry),
	})

	// 服务互联
	link.ConnectNode(&link.NodeParameter{
		PeerType: "tcp.Connector",
		NetProc:  "tcp.svc",
		SvcName:  "hub",
	})

	link.CheckReady()

	fx.WaitExit()

}
