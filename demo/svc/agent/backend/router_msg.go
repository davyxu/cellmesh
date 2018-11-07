package backend

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/agent/model"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/relay"
)

func init() {

	proto.Handle_Agent_CloseClientACK = func(ev cellnet.Event) {

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

	proto.Handle_Agent_Default = func(ev cellnet.Event) {

		//switch msg := ev.Message().(type) {
		//case *service.ServiceIdentifyACK:
		//	recoverBackend(ev.Session(), msg.SvcName)
		//case *cellnet.SessionClosed:
		//	removeBackend(ev.Session())
		//}
	}

	// 从后端服务器收到的消息
	relay.SetBroadcaster(func(ev *relay.RecvMsgEvent) {

		// 列表广播
		if value := ev.PassThroughAsInt64Slice(); value != nil {
			for _, sesid := range value {
				ses := model.GetClientSession(sesid)
				if ses != nil {
					ses.Send(ev.Msg)
				}
			}
		}

		// 原样回复
		if svcid := ev.PassThroughAsString(); svcid != "" {
			if svcid == model.AgentSvcID {
				ses := model.GetClientSession(ev.PassThroughAsInt64())
				if ses != nil {
					ses.Send(ev.Msg)
				}
			} else {
				panic(fmt.Sprintf("Recv backend msg not belong to this agent, expect '%s', got '%s'", model.AgentSvcID, svcid))
			}
			// 单发
		} else if clientSesID := ev.PassThroughAsInt64(); clientSesID != 0 {
			ses := model.GetClientSession(clientSesID)
			if ses != nil {
				ses.Send(ev.Msg)
			}
		} else {
			// TODO 只广播给认证用户?
			model.FrontendSessionManager.VisitSession(func(clientSes cellnet.Session) bool {

				clientSes.Send(ev.Message())

				return true
			})

		}

	})

}
