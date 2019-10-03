package memsd

import (
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/rpc"
)

func init() {
	// 仅供demo使用的
	proc.RegisterProcessor("memsd.cli", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcpMessageTransmitter))

		bundle.SetHooker(new(rpc.TypeRPCHooker))

		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})

	proc.RegisterProcessor("memsd.svc", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcpMessageTransmitter))

		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
