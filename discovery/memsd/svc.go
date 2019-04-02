package main

import (
	"github.com/davyxu/cellmesh/discovery/memsd/api"
	"github.com/davyxu/cellmesh/discovery/memsd/model"
	"github.com/davyxu/cellmesh/discovery/memsd/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"
)

var log = golog.New("memsd")

func startSvc() {

	config := memsd.DefaultConfig()
	if *flagAddr != "" {
		config.Address = *flagAddr
	}

	model.Queue = cellnet.NewEventQueue()
	model.Queue.EnableCapturePanic(true)
	model.Queue.StartLoop()

	p := peer.NewGenericPeer("tcp.Acceptor", "memsd", config.Address, model.Queue)
	p.(cellnet.PeerCaptureIOPanic).EnableCaptureIOPanic(true)
	model.Listener = p
	msgFunc := proto.GetMessageHandler("memsd")

	proc.BindProcessorHandler(p, "memsd.svc", func(ev cellnet.Event) {

		if msgFunc != nil {
			msgFunc(ev)
		}
	})

	p.(cellnet.TCPSocketOption).SetSocketBuffer(1024*1024, 1024*1024, true)
	p.(cellnet.PeerCaptureIOPanic).EnableCaptureIOPanic(true)
	p.Start()
	service.WaitExitSignal()
}

func deleteNotify(key, reason string) {
	valueMeta := model.DeleteValue(key)

	var ack proto.ValueDeleteNotifyACK
	ack.Key = key

	if valueMeta != nil {
		ack.SvcName = valueMeta.SvcName
	}

	if valueMeta != nil {

		if valueMeta.SvcName == "" {
			log.Infof("DeleteValue '%s'  reason: %s", key, reason)
		} else {
			log.Infof("DeregisterService '%s'  reason: %s", model.GetSvcIDByServiceKey(key), reason)
		}
	}

	model.Broadcast(&ack)

}

func checkAuth(ses cellnet.Session) bool {

	return model.GetSessionToken(ses) != ""
}
