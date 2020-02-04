package link

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/ulog"
)

// 开启服务
func ListenService(param *ServiceParameter) cellnet.Peer {

	p := peer.NewGenericPeer(param.PeerType, param.SvcName, param.ListenAddress, param.Queue)

	proc.BindProcessorHandler(p, param.NetProc, param.EventCallback)

	if opt, ok := p.(cellnet.TCPSocketOption); ok {
		opt.SetSocketBuffer(2048, 2048, true)
	}

	AddPeer(p)

	p.Start()

	registerPeerToDiscovery(p)

	return p
}

type peerListener interface {
	Port() int
}

type ServiceMeta map[string]string

// 将Acceptor注册到服务发现,IP自动取本地IP
func registerPeerToDiscovery(p cellnet.Peer, options ...interface{}) *discovery.ServiceDesc {
	host := util.GetLocalIP()

	property := p.(cellnet.PeerProperty)

	sd := &discovery.ServiceDesc{
		ID:   MakeSvcID(property.Name()),
		Name: property.Name(),
		Host: host,
		Port: p.(peerListener).Port(),
	}

	for _, opt := range options {

		switch optValue := opt.(type) {
		case ServiceMeta:
			for metaKey, metaValue := range optValue {
				sd.SetMeta(metaKey, metaValue)
			}
		}
	}

	// TODO 自动获取外网IP
	if fx.WANIP != "" {
		sd.SetMeta(SDMetaKey_WANAddress, util.JoinAddress(fx.WANIP, sd.Port))
	}

	p.(cellnet.ContextSet).SetContext(PeerContextKey_ServiceDesc, sd)

	// 有同名的要先解除注册，再注册，防止watch不触发
	discovery.Global.Deregister(sd.ID)
	err := discovery.Global.Register(sd)
	if err != nil {
		ulog.Errorf("service register failed, %s %s", sd.String(), err.Error())
	}

	return sd
}

// 解除peer注册
func unregister(p cellnet.Peer) {
	property := p.(cellnet.PeerProperty)
	discovery.Global.Deregister(MakeSvcID(property.Name()))
}
