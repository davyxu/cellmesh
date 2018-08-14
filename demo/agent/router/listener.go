package router

import (
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
)

var (
	clientListener cellnet.Peer
)

func Start() {
	clientListener = peer.NewGenericPeer("tcp.Acceptor", "", ":18000", nil)

	proc.BindProcessorHandler(clientListener, "demo.agent", nil)

	clientListener.Start()

	listenPort := clientListener.(cellnet.TCPAcceptor).ListenPort()

	name := fmt.Sprintf("agent-%d", listenPort)
	clientListener.(cellnet.PeerProperty).SetName("frontend")

	host := util.GetLocalIP()

	sd := &discovery.ServiceDesc{
		Host: host,
		Port: listenPort,
		ID:   name,
		Name: "demo.agent",
	}

	discovery.Default.Register(sd)
}

func Stop() {

	if clientListener != nil {
		clientListener.Stop()

		svcid := clientListener.(cellnet.PeerProperty).Name()

		discovery.Default.Deregister(svcid)
	}

}
