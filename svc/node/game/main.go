package main

import (
	"github.com/davyxu/cellmesh/fx/link"
	_ "github.com/davyxu/cellmesh/fx/proc"
	_ "github.com/davyxu/cellmesh/svc/node/agent/api"
)

import (
	"github.com/davyxu/cellmesh/fx"
	_ "github.com/davyxu/cellmesh/svc/node/game/enter"
)

func main() {
	fx.Init("game")
	fx.LogParameter()
	link.ConnectDiscovery()

	link.RegisterBackendNode()

	// 服务互联
	link.ConnectNode(&link.NodeParameter{
		PeerType:      "tcp.Connector",
		NetProc:       "svc.backend",
		NodeName:      "agent",
		EventCallback: fx.MakeIOCEventHandler(fx.MessageRegistry),
	})

	// 服务互联
	link.ConnectNode(&link.NodeParameter{
		PeerType: "tcp.Connector",
		NetProc:  "tcp.svc",
		NodeName: "hub",
	})

	link.CheckReady()

	fx.WaitExit()

}
