package basefx

import (
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"time"
)

var (
	// 管理Acceptor的peer，方便关闭时去掉服务发现注册
	peers []cellnet.Peer
)

func CreateCommnicateAcceptor(param fxmodel.ServiceParameter) cellnet.Peer {

	if param.NetPeerType == "" {
		param.NetPeerType = "tcp.Acceptor"
	}

	if param.NetProcName == "" {
		param.NetProcName = "tcp.svc"
	}

	p := peer.NewGenericPeer("tcp.Acceptor", param.SvcName, param.ListenAddr, fxmodel.Queue)

	msgFunc := proto.GetMessageHandler(param.SvcName)

	//"tcp.svc"
	proc.BindProcessorHandler(p, param.NetProcName, func(ev cellnet.Event) {

		if msgFunc != nil {
			msgFunc(ev)
		}
	})

	peers = append(peers, p)

	p.Start()

	service.Register(p)

	return p
}

func CreateCommnicateConnector(param fxmodel.ServiceParameter) {
	if param.NetPeerType == "" {
		param.NetPeerType = "tcp.Connector"
	}
	if param.NetProcName == "" {
		param.NetProcName = "tcp.svc"
	}

	svcName := service.GetProcName()

	msgFunc := proto.GetMessageHandler(svcName)

	opt := service.DiscoveryOption{
		MaxCount: param.MaxConnCount,
	}

	opt.Rules = service.LinkRules

	go service.DiscoveryConnector(param.SvcName, opt, func(sd *discovery.ServiceDesc) cellnet.Peer {

		p := peer.NewGenericPeer(param.NetPeerType, param.SvcName, sd.Address(), fxmodel.Queue)

		proc.BindProcessorHandler(p, param.NetProcName, func(ev cellnet.Event) {

			if msgFunc != nil {
				msgFunc(ev)
			}
		})

		p.(cellnet.TCPConnector).SetReconnectDuration(time.Second * 3)

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
