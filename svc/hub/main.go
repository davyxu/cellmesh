package main

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/link"
)

func main() {
	cellmesh.Init("hub")
	cellmesh.LogParameter()
	cellmesh.ConnectDiscovery()

	link.StartService(&link.ServiceParameter{
		ListenAddress: ":0",
	})

	link.CheckReady()

	cellmesh.WaitExitSignal()

}
