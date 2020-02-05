package main

import (
	_ "github.com/davyxu/cellmesh/fx/proc"
)

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/link"
	_ "github.com/davyxu/cellmesh/svc/agent/backend"
	_ "github.com/davyxu/cellmesh/svc/agent/frontend"
	"github.com/davyxu/cellmesh/svc/agent/routerule"
)

func main() {
	fx.Init("agent")
	fx.LogParameter()
	link.ConnectDiscovery()

	routerule.Download()

	// 网关对客户端连接
	link.ListenService(&link.ServiceParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "tcp.frontend",
		SvcName:       "frontend",
		ListenAddress: ":8002",
		Queue:         fx.Queue,
	})

	// 对内的服务连接
	link.ListenService(&link.ServiceParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "agent.backend",
		SvcName:       "backend",
		ListenAddress: ":0",
		Queue:         fx.Queue,
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
