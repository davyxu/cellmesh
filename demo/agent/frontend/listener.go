package frontend

import (
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/svcfx/model"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
)

const (
	agentSvcName = "agent"
)

func Start(add string) {
	model.FrontendListener = peer.NewGenericPeer("tcp.Acceptor", agentSvcName, add, nil)

	proc.BindProcessorHandler(model.FrontendListener, "agent.frontend", nil)

	model.FrontendListener.Start()

	listenPort := model.FrontendListener.(cellnet.TCPAcceptor).Port()

	model.FrontendListener.(cellnet.PeerProperty).SetName("frontend")

	host := util.GetLocalIP()

	sd := &discovery.ServiceDesc{
		Host: host,
		Port: listenPort,
		ID:   fxmodel.GetSvcID(agentSvcName),
		Name: agentSvcName,
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
