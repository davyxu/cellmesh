package backend

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/relay"
)

func init() {

	proto.Handle_Agent_backend_BindBackendACK = func(ev cellnet.Event) {

		msg := ev.Message().(*proto.BindBackendACK)

		bindClientToBackend(ev.Session(), msg.ID)
	}

	proto.Handle_Agent_backend_CloseClientACK = func(ev cellnet.Event) {

		msg := ev.Message().(*proto.CloseClientACK)

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

	proto.Handle_Agent_backend_Default = func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *service.ServiceIdentifyACK:
			recoverBackend(ev.Session(), msg.SvcName)
		case *cellnet.SessionClosed:
			removeBackend(ev.Session())
		}
	}

	// 从后端服务器收到的消息
	relay.SetBroadcaster(func(ev *relay.RecvMsgEvent) {

		switch tgt := ev.PassThrough.(type) {
		case int64: // 单发
			ses := model.GetClientSession(tgt)
			if ses != nil {
				ses.Send(ev.Msg)
			}
		case []int64: // 列表广播

			for _, sesid := range tgt {
				ses := model.GetClientSession(sesid)
				if ses != nil {
					ses.Send(ev.Msg)
				}
			}

		case *proto.ClientID: // 原样回复
			if tgt.SvcID == model.AgentSvcID {
				ses := model.GetClientSession(tgt.ID)
				if ses != nil {
					ses.Send(ev.Msg)
				}
			} else {
				panic(fmt.Sprintf("Recv backend msg not belong to this agent, expect '%s', got '%s'", model.AgentSvcID, tgt.SvcID))
			}
		case nil: // 广播所有

			// TODO 只广播给认证用户?
			model.FrontendSessionManager.VisitSession(func(clientSes cellnet.Session) bool {

				clientSes.Send(ev.Message())

				return true
			})

		default:
			panic("Invalid backend passthrough format")
		}

	})

}
