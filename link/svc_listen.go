package link

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/fx"
	memsd "github.com/davyxu/cellmesh/svc/memsd/api"
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
		ID:   fx.MakeSvcID(property.Name()),
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
	p.(cellnet.ContextSet).SetContext("NeedRecover", true)

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
	discovery.Global.Deregister(fx.MakeSvcID(property.Name()))
}

// 连接到服务发现, 建议在service.Init后, 以及服务器逻辑开始前调用
func ConnectDiscovery() {
	ulog.Debugf("Connecting to discovery '%s' ...", fx.DiscoveryAddress)
	sdConfig := memsd.DefaultConfig()
	sdConfig.Address = fx.DiscoveryAddress
	discovery.Global = memsd.NewDiscovery()
	discovery.Global.Start(sdConfig)

	go autoRecoverServiceDiscovery()
}

// 服务发现断线重连
func autoRecoverServiceDiscovery() {
	notify := discovery.Global.RegisterNotify()

	for {
		ctx := <-notify
		if ctx.Mode == "ready" {
			ulog.Infof("memsd discovery recover")

			VisitPeer(func(p cellnet.Peer) bool {

				// 需要重新注册的服务只有侦听器
				if _, ok := p.(cellnet.ContextSet).GetContext("NeedRecover"); ok {
					registerPeerToDiscovery(p)
				}

				return true
			})

		}

	}

}
