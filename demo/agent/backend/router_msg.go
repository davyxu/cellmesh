package backend

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
)

func init() {
	proto.Handle_Router_RouterBindUserACK = func(event service.Event, msg *proto.RouterBindUserACK) {
		bindClientToBackend(event.Session(), msg.ID)
	}

	proto.Handle_Router_Default = func(ev service.Event) {

		switch msg := ev.Message().(type) {
		case *proto.ServiceIdentifyACK:
			recoverBackend(ev.Session(), msg.SvcName)
		case *cellnet.SessionClosed:
			removeBackend(ev.Session())
		}

	}
}
