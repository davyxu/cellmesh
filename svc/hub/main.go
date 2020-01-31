package main

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/link"
)

func main() {
	cellmesh.Init("hub")
	cellmesh.LogParameter()
	cellmesh.ConnectDiscovery()

	// 服务互联服务
	link.StartService(&link.ServiceParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "tcp.svc",
		SvcName:       "hub",
		ListenAddress: ":7001",
		Queue:         cellmesh.Queue,
	})

	link.CheckReady()

	cellmesh.WaitExitSignal()

}
