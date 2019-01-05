package memsd

import (
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
)

func init() {
	// 仅供demo使用的
	proc.RegisterProcessor("memsd.cli", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(new(TCPMessageTransmitter))
		//bundle.SetHooker(proc.NewMultiHooker(new(tcp.MsgHooker), new(typeRPCHooker)))
		bundle.SetHooker(new(typeRPCHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})

	proc.RegisterProcessor("memsd.svc", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(new(TCPMessageTransmitter))
		//bundle.SetHooker(proc.NewMultiHooker(new(tcp.MsgHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
