package agentapi

import (
	"github.com/davyxu/cellmesh/link"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/ulog"
)

type AgentHooker struct {
}

// 后端服务器接收来自网关的消息
func (AgentHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch incomingMsg := inputEvent.Message().(type) {
	case *proto.RouterTransmitACK:

		userMsg, _, err := codec.DecodeMessage(int(incomingMsg.MsgID), incomingMsg.MsgData)
		if err != nil {
			ulog.Warnf("Backend msg decode failed, %s, msgid: %d", err.Error(), incomingMsg.MsgID)
			return nil
		}

		invokeAgentMessage(&AgentMsgEvent{
			Ses:      inputEvent.Session(),
			Msg:      userMsg,
			ClientID: incomingMsg.ClientID,
		})

		outputEvent = nil

	default:
		outputEvent = inputEvent
	}

	return
}

// 后端服务器发送到网关的消息
func (AgentHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch outgoingMsg := inputEvent.Message().(type) {
	case *proto.RouterTransmitACK:

		if ulog.IsLevelEnabled(ulog.DebugLevel) {
			writeAgentLog(inputEvent.Session(), "send", outgoingMsg)
		}
	}

	return inputEvent
}

func writeAgentLog(ses cellnet.Session, dir string, ack *proto.RouterTransmitACK) {

	if !msglog.IsMsgLogValid(int(ack.MsgID)) {
		return
	}

	peerInfo := ses.Peer().(cellnet.PeerProperty)

	userMsg, _, err := codec.DecodeMessage(int(ack.MsgID), ack.MsgData)
	if err == nil {
		ulog.Debugf("#agent.%s(%s)@%d len: %d %s <%d>| %s",
			dir,
			peerInfo.Name(),
			ses.ID(),
			cellnet.MessageSize(userMsg),
			cellnet.MessageToName(userMsg),
			ack.ClientID,
			cellnet.MessageToString(userMsg))
	} else {

		// 网关没有相关的消息, 只能打出消息号
		ulog.Debugf("#agent.%s(%s)@%d len: %d msgid: %d <%d>",
			dir,
			peerInfo.Name(),
			ses.ID(),
			len(ack.MsgData),
			ack.MsgID,
			ack.ClientID,
		)
	}
}
func init() {
	// 适用于后端服务的处理器
	proc.RegisterProcessor("svc.backend", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(
			new(link.SvcEventHooker), // 服务互联处理
			new(AgentHooker),         // 网关消息处理
			new(tcp.MsgHooker)))      // tcp基础消息处理
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
