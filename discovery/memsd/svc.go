package main

import (
	memsd "github.com/davyxu/cellmesh/discovery/memsd/api"
	"github.com/davyxu/cellmesh/discovery/memsd/model"
	sdproto "github.com/davyxu/cellmesh/discovery/memsd/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/golog"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
	proc.BindProcessorHandler(p, "memsd.svc", msgHandler)

	p.(cellnet.TCPSocketOption).SetSocketBuffer(1024*1024, 1024*1024, true)
	p.(cellnet.PeerCaptureIOPanic).EnableCaptureIOPanic(true)
	p.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch
}

func deleteValueRecurse(key, reason string) {

	var keyToDelete []string
	model.VisitValue(func(meta *model.ValueMeta) bool {

		if strings.HasPrefix(meta.Key, key) {
			keyToDelete = append(keyToDelete, meta.Key)
		}

		return true
	})

	for _, key := range keyToDelete {
		deleteNotify(key, reason)
	}
}

func deleteNotify(key, reason string) {
	valueMeta := model.DeleteValue(key)

	var ack sdproto.ValueDeleteNotifyACK
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
