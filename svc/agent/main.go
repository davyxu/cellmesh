package main

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/link"
	_ "github.com/davyxu/cellmesh/svc/agent/backend"
	_ "github.com/davyxu/cellmesh/svc/agent/frontend"
)

func main() {
	cellmesh.Init("agent")
	cellmesh.LogParameter()
	cellmesh.ConnectDiscovery()

	// 网关对客户端连接
	link.StartService(&link.ServiceParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "tcp.frontend",
		SvcName:       "frontend",
		ListenAddress: ":8002",
		Queue:         cellmesh.Queue,
	})

	// 对内的服务连接
	link.StartService(&link.ServiceParameter{
		PeerType:      "tcp.Acceptor",
		NetProc:       "agent.backend",
		SvcName:       "frontend",
		ListenAddress: ":0",
		Queue:         cellmesh.Queue,
	})

	// 服务互联
	link.LinkService(&link.ServiceParameter{
		PeerType: "tcp.Connector",
		NetProc:  "tcp.svc",
		SvcName:  "hub",
		Queue:    cellmesh.Queue,
	})

	link.CheckReady()

	cellmesh.WaitExitSignal()

}
