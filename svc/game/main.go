package main

import (
	_ "github.com/davyxu/cellmesh/fx/proc"
	"github.com/davyxu/cellmesh/link"
	_ "github.com/davyxu/cellmesh/svc/agent/api"
)

import (
	"github.com/davyxu/cellmesh/fx"
	_ "github.com/davyxu/cellmesh/svc/game/enter"
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
		SvcName:       "agent",
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
