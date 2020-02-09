package main

import (
	_ "github.com/davyxu/cellmesh/fx/proc"
)

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/link"
	_ "github.com/davyxu/cellmesh/svc/hub/relay"
)

func main() {
	fx.Init("hub")
	fx.LogParameter()
	link.ConnectDiscovery()

	// 服务互联服务
	link.ListenService(&link.ServiceParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "tcp.svc",
		SvcName:       "hub",
		ListenAddress: ":0",
		Queue:         fx.Queue,
		EventCallback: fx.MakeIOCEventHandler(fx.MessageRegistry),
	})

	link.CheckReady()

	fx.WaitExitSignal()

}
