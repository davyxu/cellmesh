package main

import (
	_ "github.com/davyxu/cellmesh/fx/proc"
)

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/link"
	_ "github.com/davyxu/cellmesh/svc/game/enter"
)

func main() {
	fx.Init("game")
	fx.LogParameter()
	link.ConnectDiscovery()

	// 服务互联
	link.ConnectService(&link.ServiceParameter{
		PeerType:      "tcp.Connector",
		NetProc:       "tcp.svc",
		SvcName:       "backend",
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
