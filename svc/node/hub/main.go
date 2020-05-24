package main

import (
	"github.com/davyxu/cellmesh/fx/link"
	_ "github.com/davyxu/cellmesh/fx/proc"
)

import (
	"github.com/davyxu/cellmesh/fx"
	_ "github.com/davyxu/cellmesh/svc/node/hub/relay"
)

func main() {
	fx.Init("hub")
	fx.LogParameter()
	link.ConnectDiscovery()

	// 跨服通信
	link.ListenNode(&link.NodeParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "tcp.svc",
		NodeName:      "hub",
		ListenAddress: ":0",
		EventCallback: fx.MakeIOCEventHandler(fx.MessageRegistry),
	})

	link.CheckReady()

	fx.WaitExit()

}
