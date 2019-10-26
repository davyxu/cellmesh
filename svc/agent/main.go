package main

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/link"
)

func main() {
	cellmesh.Init("agent")
	cellmesh.LogParameter()
	cellmesh.ConnectDiscovery()

	link.StartService(&link.ServiceParameter{
		ListenAddress: ":0",
	})

	link.LinkService(&link.ServiceParameter{
		SvcName: "hub",
	})

	link.CheckReady()

	cellmesh.WaitExitSignal()

}
