package main

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/link"
)

func main() {
	cellmesh.Init("login")
	cellmesh.LogParameter()
	cellmesh.ConnectDiscovery()

	link.LinkService(&link.ServiceParameter{
		SvcName: "hub",
	})

	link.CheckReady()

	cellmesh.WaitExitSignal()

}
