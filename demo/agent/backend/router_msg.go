package backend

import (
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
)

func init() {

	proto.Handle_Router_BindBackendACK = func(event service.Event, msg *proto.BindBackendACK) {
		bindClientToBackend(event.Session(), msg.ID)
	}

	proto.Handle_Router_CloseClientACK = func(ev service.Event, msg *proto.CloseClientACK) {

		// 不给ID,掐线这个网关的所有客户端
		if len(msg.ID) == 0 {
			model.VisitUser(func(user *model.User) bool {
				user.ClientSession.Close()
				return true
			})

		} else {
			// 关闭指定的客户端
			for _, sesid := range msg.ID {
				if u := model.GetUser(sesid); u != nil {
					u.ClientSession.Close()
				}
			}

		}

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
