package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"sync"
)

// 服务互联消息处理
type SvcEventHooker struct {
}

func (SvcEventHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch msg := inputEvent.Message().(type) {
	case *ServiceIdentifyACK:

		if pre := GetRemoteService(msg.SvcID); pre == nil {

			// 添加连接上来的对方服务
			AddRemoteService(inputEvent.Session(), &discovery.ServiceDesc{
				ID:   msg.SvcID,
				Name: msg.SvcName,
			})
		}
	case *cellnet.SessionConnected:

		ctx := inputEvent.Session().Peer().(cellnet.ContextSet)

		var sd *discovery.ServiceDesc
		if !ctx.FetchContext("sd", &sd) {
			panic("sd get failed")
		}

		property := inputEvent.Session().Peer().(cellnet.PeerProperty)

		// 用Connector的名称（一般是ProcName）让远程知道自己是什么服务，用于网关等需要反向发送消息的标识
		inputEvent.Session().Send(ServiceIdentifyACK{
			SvcName: property.Name(),
			SvcID:   MakeServiceID(property.Name()),
		})

		AddRemoteService(inputEvent.Session(), sd)

	case *cellnet.SessionClosed:
		ctx := inputEvent.Session().Peer().(cellnet.ContextSet)

		// 只有连接器做这个
		var connMap *sync.Map
		if ctx.FetchContext("connMap", &connMap) {

			var sd *discovery.ServiceDesc
			if ctx.FetchContext("sd", &sd) {
				connMap.Delete(sd.ID)
			} else {
				panic("sd get failed")
			}
		}

		RemoveRemoteService(inputEvent.Session())
	}

	return inputEvent

}

func (SvcEventHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}

func init() {

	transmitter := new(tcp.TCPMessageTransmitter)
	svcHooker := new(SvcEventHooker)
	msgHooker := new(tcp.MsgHooker)

	proc.RegisterProcessor("tcp.svc", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(transmitter)
		bundle.SetHooker(proc.NewMultiHooker(svcHooker, msgHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
