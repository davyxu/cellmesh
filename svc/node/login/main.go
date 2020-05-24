package main

import (
	"github.com/davyxu/cellmesh/fx/link"
	_ "github.com/davyxu/cellmesh/fx/proc"
	"time"
)

import (
	"github.com/davyxu/cellmesh/fx"
	_ "github.com/davyxu/cellmesh/svc/node/login/verify"
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
		NodeName:      "verify",
		ListenAddress: ":8001",
		EventCallback: fx.MakeIOCEventHandler(fx.MessageRegistry),
	})

	// 跨服通信
	link.ConnectNode(&link.NodeParameter{
		PeerType: "tcp.Connector",
		NetProc:  "tcp.svc",
		NodeName: "hub",
	})

	link.CheckReady()

	fx.WaitExit()

}
