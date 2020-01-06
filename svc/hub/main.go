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
		SvcName:       "hub",
		ListenAddress: ":0",
	})

	link.CheckReady()

	cellmesh.WaitExitSignal()

}
