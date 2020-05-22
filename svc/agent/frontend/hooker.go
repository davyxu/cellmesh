package frontend

import (
	"fmt"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellmesh/svc/agent/model"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/ulog"

	"github.com/davyxu/cellnet/proc"
	"time"
)

var (
	PingACKMsgID        = cellnet.MessageMetaByFullName("proto.PingACK").ID
	BindBackendREQMsgID = cellnet.MessageMetaByFullName("proto.AgentBindBackendREQ").ID
)

func ProcFrontendPacket(msgID int, msgData []byte, ses cellnet.Session) (msg interface{}, err error) {
	// agent自己的内部消息以及预处理消息
	switch int(msgID) {
	case PingACKMsgID, BindBackendREQMsgID:

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
		case *proto.AgentBindBackendREQ:
			model.SendToBackend(userMsg.NodeID, msgID, msgData, ses.ID())
		}

	default:
		// 在路由规则中查找消息ID是否是路由规则允许的消息
		rule := model.GetRuleByMsgID(msgID)
		if rule == nil {
			return nil, fmt.Errorf("Message not in route table, msgid: %d, execute MakeProto.sh and restart agent", msgID)
		}

		// 找到已经绑定的用户
		u := model.SessionToUser(ses)

		if u != nil {
			// 透传到后台
			u.SendToBackend(u.GetBackend(rule.SvcName), msgID, msgData)
		} else {
			ulog.Warnf("User not bind to backend,  msgname: %s", msgID, rule.MsgName)
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
			// TODO 后端服务向网关订阅客户端断开通知, 否则不通知
			u.BroadcastToBackends(&proto.AgentClientClosedNotifyACK{
				ID: proto.AgentClientID{
					SessionID: inputEvent.Session().ID(),
					NodeID:    model.AgentNodeID,
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
	proc.RegisterProcessor("tcp.frontend", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(directTCPTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(
			new(FrontendEventHooker), // 内部消息处理
		))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
