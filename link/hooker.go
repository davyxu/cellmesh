package link

import (
	"github.com/davyxu/cellmesh/fx"
	meshproto "github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
)

// 服务互联消息处理
type SvcEventHooker struct {
}

func (SvcEventHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch msg := inputEvent.Message().(type) {
	case *meshproto.ServiceIdentifyACK: // 服务方收到连接方的服务标识

		if pre := GetLink(msg.SvcID); pre == nil {

			// 添加连接上来的对方服务
			markLink(inputEvent.Session(), msg.SvcID, msg.SvcName)
		}
	case *cellnet.SessionConnected:

		// 用Connector的名称（一般是ProcName）让远程知道自己是什么服务，用于网关等需要反向发送消息的标识
		inputEvent.Session().Send(&meshproto.ServiceIdentifyACK{
			SvcID:   fx.LocalSvcID,
			SvcName: fx.ProcName,
		})

	case *cellnet.SessionClosed:
	}

	return inputEvent

}

func (SvcEventHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}
