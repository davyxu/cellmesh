package main

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/linkmgr"
)

func main() {
	cellmesh.Init("hub")
	cellmesh.LogParameter()
	cellmesh.ConnectDiscovery()

	linkmgr.StartService(linkmgr.ServiceParameter{
		ListenAddress: ":0",
	})

	linkmgr.CheckReady()

	cellmesh.WaitExitSignal()

}
