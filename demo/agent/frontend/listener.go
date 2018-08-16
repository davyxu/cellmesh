package frontend

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
)

func Start() {
	model.FrontendListener = peer.NewGenericPeer("tcp.Acceptor", "demo.agent", ":18000", nil)

	proc.BindProcessorHandler(model.FrontendListener, "demo.agent", nil)

	model.FrontendListener.Start()

	listenPort := model.FrontendListener.(cellnet.TCPAcceptor).Port()

	name := fmt.Sprintf("agent-%d", listenPort)
	model.FrontendListener.(cellnet.PeerProperty).SetName("frontend")

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

	if model.FrontendListener != nil {
		model.FrontendListener.Stop()

		svcid := model.FrontendListener.(cellnet.PeerProperty).Name()

		discovery.Default.Deregister(svcid)
	}

}
