package main

import (
	_ "github.com/davyxu/cellmesh/fx/proc"
	"github.com/davyxu/cellmesh/link"
	"github.com/davyxu/cellmesh/svc/agent/model"
	"github.com/davyxu/cellnet/peer"
)

import (
	"github.com/davyxu/cellmesh/fx"
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
	model.FrontendSessionManager = link.ListenNode(&link.NodeParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "tcp.frontend",
		SvcName:       "frontend",
		ListenAddress: ":8002",
		Queue:         fx.Queue,
	}).(peer.SessionManager)

	// 对内的服务连接
	link.ListenNode(&link.NodeParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "agent.backend",
		SvcName:       "backend",
		ListenAddress: ":0",
		Queue:         fx.Queue,
	})

	// 服务互联
	link.ConnectNode(&link.NodeParameter{
		PeerType: "tcp.Connector",
		NetProc:  "tcp.svc",
		SvcName:  "hub",
		Queue:    fx.Queue,
	})

	link.CheckReady()

	fx.WaitExit()

}
