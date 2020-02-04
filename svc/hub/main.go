package main

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/link"
)

func main() {
	fx.Init("hub")
	fx.LogParameter()
	fx.ConnectDiscovery()

	// 服务互联服务
	link.ListenService(&link.ServiceParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "tcp.svc",
		SvcName:       "hub",
		ListenAddress: ":0",
		Queue:         fx.Queue,
	})

	link.CheckReady()

	fx.WaitExitSignal()

}
