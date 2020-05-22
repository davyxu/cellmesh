package link

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellmesh/redsd"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/ulog"
)

// 开启服务
func ListenNode(param *NodeParameter) *redsd.NodeDesc {

	p := peer.NewGenericPeer(param.PeerType, param.SvcName, param.ListenAddress, param.Queue)

	proc.BindProcessorHandler(p, param.NetProc, param.EventCallback)

	if opt, ok := p.(cellnet.TCPSocketOption); ok {
		opt.SetSocketBuffer(2048, 2048, true)
	}

	p.Start()

	return registerPeerToDiscovery(p)
}

type peerListener interface {
	Port() int
}

type ServiceMeta map[string]string

// 将Acceptor注册到服务发现,IP自动取本地IP
func registerPeerToDiscovery(p cellnet.Peer, options ...interface{}) *redsd.NodeDesc {

	property := p.(cellnet.PeerProperty)

	sd := redsd.NewDesc()
	sd.ID = fx.MakeSvcID(property.Name())
	sd.Name = property.Name()
	sd.Host = util.GetLocalIP()
	sd.Port = p.(peerListener).Port()
	sd.Peer = p

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
		sd.SetMeta("WAN", util.JoinAddress(fx.WANIP, sd.Port))
	}

	SD.NewNodeList(sd.Name, int(proto.NodeKind_Listen)).Register(sd)
	return sd
}

// 不侦听的服务, 单独注册
func RegisterBackendNode() {
	sd := redsd.NewDesc()
	sd.ID = fx.MakeSvcID(fx.ProcName)
	sd.Name = fx.ProcName
	SD.NewNodeList(sd.Name, int(proto.NodeKind_Backend)).Register(sd)
}

// 连接到服务发现, 建议在service.Init后, 以及服务器逻辑开始前调用
func ConnectDiscovery() {
	ulog.Debugf("Connecting to discovery '%s' ...", fx.DiscoveryAddress)
	SD = redsd.NewRedisDiscovery()
	SD.Start(fx.DiscoveryAddress)

	//go autoRecoverServiceDiscovery()
}

//// 服务发现断线重连
//func autoRecoverServiceDiscovery() {
//	notify := discovery.Global.RegisterNotify()
//
//	for {
//		ctx := <-notify
//		if ctx.Mode == "ready" {
//			ulog.Infof("memsd discovery recover")
//
//			VisitPeer(func(p cellnet.Peer) bool {
//
//				// 需要重新注册的服务只有侦听器
//				if _, ok := p.(cellnet.ContextSet).GetContext("NeedRecover"); ok {
//					registerPeerToDiscovery(p)
//				}
//
//				return true
//			})
//
//		}
//
//	}
//
//}
