package basefx

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
)

var (
	// 管理Acceptor的peer，方便关闭时去掉服务发现注册
	peers []cellnet.Peer
)

func CreateCommnicateAcceptor(svcName, listenAddr string) cellnet.Peer {

	p := peer.NewGenericPeer("tcp.Acceptor", svcName, listenAddr, nil)

	msgFunc := proto.GetMessageHandler(svcName)

	proc.BindProcessorHandler(p, "tcp.svc", func(ev cellnet.Event) {

		if msgFunc != nil {
			msgFunc(ev)
		}
	})

	p.Start()

	service.Register(p)
	peers = append(peers, p)
	return p
}

func CreateCommnicateConnector(tgtSvcName string) {

	svcName := service.GetProcName()

	msgFunc := proto.GetMessageHandler(svcName)

	go service.DiscoveryConnector(service.LinkRules, tgtSvcName, -1, func(sd *discovery.ServiceDesc) cellnet.Peer {

		p := peer.NewGenericPeer("tcp.Connector", svcName, sd.Address(), nil)

		proc.BindProcessorHandler(p, "tcp.svc", func(ev cellnet.Event) {

			if msgFunc != nil {
				msgFunc(ev)
			}
		})

		p.Start()

		return p
	})

}

func StopAllPeers() {
	for _, p := range peers {
		service.Unregister(p)
		p.Stop()
	}
}
