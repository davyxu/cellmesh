package main

import (
	"github.com/davyxu/cellmesh/fx/db"
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

	db.Redis.Connect()

	// 网关
	link.ConnectNode(&link.NodeParameter{
		PeerType:      "tcp.Connector",
		NetProc:       "svc.backend",
		NodeName:      "agent",
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
