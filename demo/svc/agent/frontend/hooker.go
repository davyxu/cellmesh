package frontend

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/agent/model"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"time"
)

var (
	PingACKMsgID   = cellnet.MessageMetaByFullName("proto.PingACK").ID
	VerifyREQMsgID = cellnet.MessageMetaByFullName("proto.VerifyREQ").ID
)

func ProcFrontendPacket(msgID int, msgData []byte, ses cellnet.Session) (msg interface{}, err error) {
	// agent自己的内部消息以及预处理消息
	switch int(msgID) {
	case PingACKMsgID, VerifyREQMsgID:

		// 将字节数组和消息ID用户解出消息
		msg, _, err = codec.DecodeMessage(msgID, msgData)
		if err != nil {
			// TODO 接收错误时，返回消息
			return nil, err
		}

		switch userMsg := msg.(type) {
		case *proto.PingACK:
			u := model.SessionToUser(ses)
			if u != nil {
				u.LastPingTime = time.Now()

				// 回消息
				ses.Send(&proto.PingACK{})
			} else {
				ses.Close()
			}

			// 第一个到网关的消息
		case *proto.VerifyREQ:
			u, err := bindClientToBackend(userMsg.GameSvcID, ses.ID())
			if err == nil {
				u.TransmitToBackend(userMsg.GameSvcID, msgID, msgData)

			} else {
				ses.Close()
				log.Errorln("bindClientToBackend", err)
			}
		}

	default:
		// 在路由规则中查找消息ID是否是路由规则允许的消息
		rule := model.GetRuleByMsgID(msgID)
		if rule == nil {
			return nil, fmt.Errorf("Message not in route table, msgid: %d, use MakeProto.sh!", msgID)
		}

		// 找到已经绑定的用户
		u := model.SessionToUser(ses)

		if u != nil {

			// 透传到后台
			if err = u.TransmitToBackend(u.GetBackend(rule.SvcName), msgID, msgData); err != nil {
				log.Warnf("TransmitToBackend %s, msg: '%s' svc: %s", err, rule.MsgName, rule.SvcName)
			}

		} else {
			// 这是一个未授权的用户发授权消息,可以踢掉
		}
	}

	return
}

type FrontendEventHooker struct {
}

// 网关内部抛出的事件
func (FrontendEventHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch inputEvent.Message().(type) {
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
	}

	return inputEvent
}

// 发送到客户端的消息
func (FrontendEventHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}

func init() {

	// 前端的processor
	proc.RegisterProcessor("agent.frontend", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(new(directTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(
			new(tcp.MsgHooker),       //  TCP基础消息及日志
			new(FrontendEventHooker), // 内部消息处理
		))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
