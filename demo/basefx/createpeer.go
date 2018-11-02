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

func CreateCommnicateAcceptor(param fxmodel.ServiceParameter) cellnet.Peer {

	if param.NetPeerType == "" {
		param.NetPeerType = "tcp.Acceptor"
	}

	if param.NetProcName == "" {
		param.NetProcName = "tcp.svc"
	}

	var q cellnet.EventQueue
	if !param.NoQueue {
		q = fxmodel.Queue
	}

	p := peer.NewGenericPeer("tcp.Acceptor", param.SvcName, param.ListenAddr, q)

	msgFunc := proto.GetMessageHandler(param.SvcName)

	//"tcp.svc"
	proc.BindProcessorHandler(p, param.NetProcName, func(ev cellnet.Event) {

		if msgFunc != nil {
			msgFunc(ev)
		}
	})

	if opt, ok := p.(cellnet.TCPSocketOption); ok {
		opt.SetSocketBuffer(2048, 2048, true)
	}

	fxmodel.AddLocalService(p)

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

	msgFunc := proto.GetMessageHandler(service.GetProcName())

	opt := service.DiscoveryOption{
		MaxCount: param.MaxConnCount,
	}

	// 强制匹配自己的组
	if *fxmodel.FlagSelfGroup {
		opt.MatchSvcGroup = service.GetSvcGroup()
	}

	opt.Rules = service.LinkRules

	var q cellnet.EventQueue
	if !param.NoQueue {
		q = fxmodel.Queue
	}

	mp := service.DiscoveryService(param.SvcName, opt, func(sd *discovery.ServiceDesc) cellnet.Peer {

		p := peer.NewGenericPeer(param.NetPeerType, param.SvcName, sd.Address(), q)

		proc.BindProcessorHandler(p, param.NetProcName, func(ev cellnet.Event) {

			if msgFunc != nil {
				msgFunc(ev)
			}
		})

		if opt, ok := p.(cellnet.TCPSocketOption); ok {
			opt.SetSocketBuffer(2048, 2048, true)
		}

		p.(cellnet.TCPConnector).SetReconnectDuration(time.Second * 3)

		p.Start()

		return p
	})

	mp.(service.MultiPeer).SetContext(param)

	fxmodel.AddLocalService(mp)

}
