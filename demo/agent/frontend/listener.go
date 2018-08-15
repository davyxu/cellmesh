package frontend

import (
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
)

var (
	frontendListener cellnet.Peer
)

func Start() {
	frontendListener = peer.NewGenericPeer("tcp.Acceptor", "demo.agent", ":18000", nil)

	proc.BindProcessorHandler(frontendListener, "demo.agent", nil)

	frontendListener.Start()

	listenPort := frontendListener.(cellnet.TCPAcceptor).ListenPort()

	name := fmt.Sprintf("agent-%d", listenPort)
	frontendListener.(cellnet.PeerProperty).SetName("frontend")

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

	if frontendListener != nil {
		frontendListener.Stop()

		svcid := frontendListener.(cellnet.PeerProperty).Name()

		discovery.Default.Deregister(svcid)
	}

}
