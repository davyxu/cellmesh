package link

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
)

type eventHandler interface {
	OnEvent(inputEvent cellnet.Event)
}

type MsgHooker struct {
	evHandler eventHandler
}

func (self *MsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	self.evHandler.OnEvent(inputEvent)

	msglog.WriteRecvLogger("tcp", inputEvent.Session(), inputEvent.Message())

	return inputEvent
}

func (self *MsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	msglog.WriteSendLogger("tcp", inputEvent.Session(), inputEvent.Message())

	return inputEvent
}

func NewMsgHooker(evHandler eventHandler) *MsgHooker {
	return &MsgHooker{
		evHandler: evHandler,
	}
}

func init() {
	// 仅供demo使用的
	proc.RegisterProcessor("tcp.robot", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(
			NewMsgHooker(args[0].(eventHandler)),
		))

	})
}
