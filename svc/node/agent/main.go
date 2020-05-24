package main

import (
	"github.com/davyxu/cellmesh/fx/link"
	_ "github.com/davyxu/cellmesh/fx/proc"
	"github.com/davyxu/cellmesh/svc/node/agent/model"
	"github.com/davyxu/cellmesh/svc/node/agent/routerule"
	"github.com/davyxu/cellnet/peer"
)

import (
	"github.com/davyxu/cellmesh/fx"
	_ "github.com/davyxu/cellmesh/svc/node/agent/backend"
	_ "github.com/davyxu/cellmesh/svc/node/agent/frontend"
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
		NodeName:      "frontend",
		ListenAddress: ":8002",
	}).Peer.(peer.SessionManager)

	// 对内的服务连接
	model.AgentNodeID = link.ListenNode(&link.NodeParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "agent.backend",
		NodeName:      "agent",
		ListenAddress: ":0",
	}).ID

	// 服务互联
	link.ConnectNode(&link.NodeParameter{
		PeerType: "tcp.Connector",
		NetProc:  "tcp.svc",
		NodeName: "hub",
	})

	link.CheckReady()

	fx.WaitExit()

}
