package hubapi

import (
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/demo/svc/hub/model"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
)

type subscriberHooker struct {
}

func (self subscriberHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch inputEvent.Message().(type) {
	case *cellnet.SessionConnected: // 连接上hub时

		model.HubSession = inputEvent.Session()

		// 自动订阅进程名和svcid对应的频道
		Subscribe(service.GetProcName())
		Subscribe(service.GetLocalSvcID())

		// 玩家自定义频道
		if model.OnHubReady != nil {
			cellnet.QueuedCall(fxmodel.Queue, model.OnHubReady)
		}

	}

	return inputEvent
}

func (self subscriberHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}

func init() {

	proc.RegisterProcessor("tcp.hub", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(subscriberHooker), new(tcp.MsgHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
