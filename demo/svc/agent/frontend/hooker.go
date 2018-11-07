package frontend

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/agent/model"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"reflect"
	"time"
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
		u := model.SessionToUser(inputEvent.Session())
		if u != nil {
			u.BroadcastToBackends(&proto.ClientClosedACK{
				ID: proto.ClientID{
					ID:    inputEvent.Session().ID(),
					SvcID: model.AgentSvcID,
				},
			})
		}
	case *proto.PingACK:
		u := model.SessionToUser(inputEvent.Session())
		if u != nil {
			u.LastPingTime = time.Now()

			// 回消息
			inputEvent.Session().Send(&proto.PingACK{})
		}
	case *proto.VerifyREQ:

		u, err := bindClientToBackend(incomingMsg.GameSvcID, inputEvent.Session().ID())
		if err == nil {

			u.RelayToService(incomingMsg.GameSvcID, incomingMsg)

		} else {
			log.Errorln("bindClientToBackend", err)
		}

	default:
		msgType := reflect.TypeOf(incomingMsg).Elem()

		// 确定消息所在的服务
		if rule := model.GetTargetService(msgType.Name()); rule != nil {

			switch rule.Mode {
			case "auth":

				// 从客户端过来的会话取得绑定的用户
				u := model.SessionToUser(inputEvent.Session())

				if u != nil {

					u.RelayToService(u.GetBackend(rule.SvcName), incomingMsg)

				} else {
					// 这是一个未授权的用户发授权消息,可以踢掉
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

	transmitter := new(tcp.TCPMessageTransmitter)
	routerHooker := new(RelayUpMsgHooker)
	msgLogger := new(tcp.MsgHooker)

	// 前端的processor
	proc.RegisterProcessor("tcp.frontend", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(transmitter)
		bundle.SetHooker(proc.NewMultiHooker(msgLogger, routerHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
