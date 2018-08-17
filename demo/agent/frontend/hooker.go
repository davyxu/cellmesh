package frontend

import (
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/cellnet/relay"
	"reflect"
)

type RelayUpMsgHooker struct {
}

/* 网关通过规则
1. 直接接收，根据消息选择后台服务地址，适用于未绑定用户的消息
2. 绑定消息，直接获得用户绑定的后台地址转发


消息名-> 服务(进程) -> Session
*/

// 从客户端接收到的消息
func (RelayUpMsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch incomingMsg := inputEvent.Message().(type) {
	case *cellnet.SessionAccepted:
	case *cellnet.SessionClosed:

		// 通知后台客户端关闭
		u := model.GetUser(inputEvent.Session())
		if u != nil {
			for _, backend := range u.Targets {
				backend.Session.Send(proto.ClientClosedACK{
					ID: inputEvent.Session().ID(),
				})
			}
		}

	default:
		msgType := reflect.TypeOf(incomingMsg).Elem()

		// 确定消息所在的服务
		if rule := model.GetTargetService(msgType.Name()); rule != nil {

			switch rule.Mode {
			case "pass":
				// TODO 挑选一台
				service.VisitConn(func(ses cellnet.Session, desc *discovery.ServiceDesc) {

					if desc.Name == rule.SvcName {
						// 透传消息
						relay.Relay(ses, incomingMsg, inputEvent.Session().ID())
					}
				})

			case "auth":

				u := model.GetUser(inputEvent.Session())

				if u != nil {

					backendSes := u.GetBackend(rule.SvcName)

					if backendSes != nil {
						relay.Relay(backendSes, incomingMsg, inputEvent.Session().ID())
					} else {
						log.Warnf("Route target not found, msg: '%s' mode: 'auth'", msgType.Name())
					}

				}

			}

		} else {

			log.Warnf("Route target not found: %s", msgType.Name())
		}
	}

	return inputEvent
}

// 发送到客户端的消息
func (RelayUpMsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}

func init() {

	// 从后端服务器收到的消息
	relay.SetBroadcaster(func(event *relay.RecvMsgEvent) {

		for _, sesID := range event.ContextID {

			ses := model.FrontendListener.(peer.SessionManager).GetSession(sesID)
			if ses == nil {
				continue
			}

			ses.Send(event.Msg)
		}

	})

	transmitter := new(tcp.TCPMessageTransmitter)
	routerHooker := new(RelayUpMsgHooker)
	msgLogger := new(tcp.MsgHooker)

	// 前端的processor
	proc.RegisterProcessor("agent.frontend", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(transmitter)
		bundle.SetHooker(proc.NewMultiHooker(msgLogger, routerHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
