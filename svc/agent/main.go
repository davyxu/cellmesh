package main

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/link"
)

func main() {
	cellmesh.Init("agent")
	cellmesh.LogParameter()
	cellmesh.ConnectDiscovery()

	// 网关对客户端连接
	link.StartService(&link.ServiceParameter{
		SvcName:       "frontend",
		ListenAddress: ":0",
	})

	// 服务互联
	link.LinkService(&link.ServiceParameter{
		SvcName: "hub",
	})

	link.CheckReady()

	cellmesh.WaitExitSignal()

}
