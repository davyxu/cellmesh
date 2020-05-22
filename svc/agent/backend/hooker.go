package backend

import (
	"github.com/davyxu/cellmesh/link"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellmesh/svc/agent/model"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/ulog"
)

var (
	BindBackendACKMsgID = cellnet.MessageMetaByFullName("proto.AgentBindBackendACK").ID
)

type backendHooker struct {
}

// 来自后台服务器的消息
func (backendHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch incomingMsg := inputEvent.Message().(type) {
	case *proto.AgentTransmitACK:

		rawPkt := &cellnet.RawPacket{
			MsgData: incomingMsg.MsgData,
			MsgID:   int(incomingMsg.MsgID),
		}

		// 后端服务打通路由通道
		if int(incomingMsg.MsgID) == BindBackendACKMsgID {
			rawMsg, _, err := codec.DecodeMessage(int(incomingMsg.MsgID), incomingMsg.MsgData)
			if err != nil {
				ulog.Errorf("decode BindBackendACK failed, %s", err)
				return nil
			}

			bindMsg := rawMsg.(*proto.AgentBindBackendACK)
			if bindMsg.Code == 0 {
				bindClientToBackend(bindMsg.NodeID, incomingMsg.ClientID)
			}
		}

		if ulog.IsLevelEnabled(ulog.DebugLevel) {
			writeBackendLog(inputEvent.Session(), "recv", incomingMsg)
		}

		// 单发
		if incomingMsg.ClientID != 0 {
			clientSes := model.GetClientSession(incomingMsg.ClientID)

			if clientSes != nil {
				clientSes.Send(rawPkt)
			}
			// 广播
		} else if incomingMsg.ClientIDList != nil {

			for _, cid := range incomingMsg.ClientIDList {
				clientSes := model.GetClientSession(cid)

				if clientSes != nil {
					clientSes.Send(rawPkt)
				}
			}
		} else if incomingMsg.All {
			model.FrontendSessionManager.VisitSession(func(clientSes cellnet.Session) bool {

				clientSes.Send(rawPkt)
				return true
			})
		}

		// 本事件已经处理, 不再后传
		return nil
	default:
		msglog.WriteRecvLogger("tcp", inputEvent.Session(), inputEvent.Message())
	}

	return inputEvent
}

// 发送给后台服务器
func (backendHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch outgoingMsg := inputEvent.Message().(type) {
	case *proto.AgentTransmitACK:

		if ulog.IsLevelEnabled(ulog.DebugLevel) {
			writeBackendLog(inputEvent.Session(), "send", outgoingMsg)
		}
	default:
		msglog.WriteSendLogger("tcp", inputEvent.Session(), inputEvent.Message())
	}

	return inputEvent
}

func writeBackendLog(ses cellnet.Session, dir string, ack *proto.AgentTransmitACK) {

	if !msglog.IsMsgLogValid(int(ack.MsgID)) {
		return
	}

	var backendNodeID string
	desc := link.DescByLink(ses)
	if desc != nil {
		backendNodeID = desc.ID
	}

	userMsg, _, err := codec.DecodeMessage(int(ack.MsgID), ack.MsgData)
	if err == nil {
		ulog.Debugf("#%s(%s) len: %d %s cid:%d| %s",
			dir,
			backendNodeID,
			cellnet.MessageSize(userMsg),
			cellnet.MessageToName(userMsg),
			ack.ClientID,
			cellnet.MessageToString(userMsg))
	} else {

		// 网关没有相关的消息, 只能打出消息号
		ulog.Debugf("#%s(%s) len: %d msgid: %d cid:%d",
			dir,
			backendNodeID,
			len(ack.MsgData),
			ack.MsgID,
			ack.ClientID,
		)
	}
}

func init() {

	// 避免默认消息日志显示本条消息
	msglog.SetMsgLogRule("proto.TransmitACK", "black")

	// 适用于
	proc.RegisterProcessor("agent.backend", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(
			new(link.SvcEventHooker), // 服务互联处理
			new(backendHooker),       // 网关消息处理
		))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
