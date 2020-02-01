package main

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/cellnet/rpc"
	"sync"
)

type connector interface {
	cellnet.TCPConnector
	cellnet.PeerReadyChecker
}

// 建立短连接
func CreateConnection(serviceName, netPeerType, netProcName, address string) (ret cellnet.Session) {

	p := peer.NewGenericPeer(netPeerType, serviceName, address, nil)
	proc.BindProcessorHandler(p, netProcName, nil)

	p.Start()

	conn := p.(connector)

	if conn.IsReady() {
		ret = conn.Session()

		return
	}

	p.Stop()
	return
}

// 保持长连接
func KeepConnection(svcid, addr, netPeerType, netProc string, onReady func(cellnet.Session), onClose func()) {

	var stop sync.WaitGroup

	p := peer.NewGenericPeer(netPeerType, svcid, addr, nil)
	proc.BindProcessorHandler(p, netProc, func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionClosed:
			stop.Done()
		}
	})

	stop.Add(1)

	p.Start()

	conn := p.(connector)

	if conn.IsReady() {

		onReady(conn.Session())

		// 连接断开
		stop.Wait()
	}

	p.Stop()

	if onClose != nil {
		onClose()
	}

}

func init() {
	// 仅供demo使用的
	proc.RegisterProcessor("tcp.demo", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(tcp.MsgHooker), new(rpc.TypeRPCHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
