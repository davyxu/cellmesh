package link

import (
	"github.com/davyxu/cellmesh/fx"
	meshproto "github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellmesh/redsd"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/ulog"
)

// 服务互联消息处理
type SvcEventHooker struct {
}

func (SvcEventHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch msg := inputEvent.Message().(type) {
	case *meshproto.ServiceIdentifyACK: // 服务方收到连接方的服务标识

		pre := DescByID(msg.SvcID)
		if pre != nil {
			ulog.Debugf("discard pre node: %s", msg.SvcID)
			closeNode(pre)
			removeLink(pre)
		}

		var desc redsd.NodeDesc
		desc.Name = msg.SvcName
		desc.ID = msg.SvcID
		desc.Session = inputEvent.Session()

		addr, _ := util.GetRemoteAddrss(inputEvent.Session())
		addobj, err := util.ParseAddress(addr)
		if err == nil {
			desc.Host = addobj.Host
			desc.Port = addobj.MinPort
		}

		ctxSet := inputEvent.Session().(cellnet.ContextSet)
		ctxSet.SetContext("NodeDesc", &desc)

		ulog.WithField("nodeid", desc.ID).Debugf("accept add link")
		addLink(&desc)

	case *cellnet.SessionConnected:

		// 从Peer上转到Session上绑定
		if raw, ok := inputEvent.Session().Peer().(cellnet.ContextSet).GetContext("NodeDesc"); ok {
			ctxSet := inputEvent.Session().(cellnet.ContextSet)
			ctxSet.SetContext("NodeDesc", raw)
		}

		// 用Connector的名称（一般是ProcName）让远程知道自己是什么服务，用于网关等需要反向发送消息的标识
		inputEvent.Session().Send(&meshproto.ServiceIdentifyACK{
			SvcID:   fx.LocalSvcID,
			SvcName: fx.ProcName,
		})

	case *cellnet.SessionClosed:

		// 这里只处理Acceptor的情况, 而Connector考虑到自动重连情况, 只由服务发现处理
		if inputEvent.Session().Peer().TypeName() == "tcp.Acceptor" {
			if raw, ok := inputEvent.Session().(cellnet.ContextSet).GetContext("NodeDesc"); ok {
				desc := raw.(*redsd.NodeDesc)
				ulog.WithField("nodeid", desc.ID).Debugf("accept remove link")
				desc.Session = nil
				removeLink(desc)
			}
		}

	}

	return inputEvent

}

func (SvcEventHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}
