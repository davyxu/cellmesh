package main

import (
	_ "github.com/davyxu/cellmesh/fx/proc"
	"github.com/davyxu/cellmesh/link"
)

import (
	"github.com/davyxu/cellmesh/fx"
	_ "github.com/davyxu/cellmesh/svc/hub/relay"
)

func main() {
	fx.Init("hub")
	fx.LogParameter()
	link.ConnectDiscovery()

	// 服务互联服务
	link.ListenNode(&link.NodeParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "tcp.svc",
		SvcName:       "hub",
		ListenAddress: ":0",
		Queue:         fx.Queue,
		EventCallback: fx.MakeIOCEventHandler(fx.MessageRegistry),
	})

	link.CheckReady()

	fx.WaitExit()

}
