package main

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/linkmgr"
)

func main() {
	cellmesh.Init("agent")
	cellmesh.LogParameter()
	cellmesh.ConnectDiscovery()

	linkmgr.StartService(linkmgr.ServiceParameter{
		ListenAddress: ":0",
	})

	linkmgr.LinkService(linkmgr.ServiceParameter{
		SvcName: "hub",
	})

	linkmgr.CheckReady()

	cellmesh.WaitExitSignal()

}
