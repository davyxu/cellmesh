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

		if ulog.IsLevelEnabled(ulog.DebugLevel) {
			writeAgentLog(inputEvent.Session(), "recv", incomingMsg)
		}

		invokeAgentMessage(&AgentMsgEvent{
			Ses:      inputEvent.Session(),
			Msg:      userMsg,
			ClientID: incomingMsg.ClientID,
		})

		outputEvent = nil

	default:

		msglog.WriteRecvLogger("tcp", inputEvent.Session(), inputEvent.Message())

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
	default:
		msglog.WriteSendLogger("tcp", inputEvent.Session(), inputEvent.Message())
	}

	return inputEvent
}

func writeAgentLog(ses cellnet.Session, dir string, ack *proto.RouterTransmitACK) {

	if !msglog.IsMsgLogValid(int(ack.MsgID)) {
		return
	}

	var agentNodeID string
	agentDesc := link.DescByLink(ses)
	if agentDesc != nil {
		agentNodeID = agentDesc.ID
	}

	userMsg, _, err := codec.DecodeMessage(int(ack.MsgID), ack.MsgData)
	if err == nil {
		ulog.Debugf("#%s(%s) len: %d %s cid:%d| %s",
			dir,
			agentNodeID,
			cellnet.MessageSize(userMsg),
			cellnet.MessageToName(userMsg),
			ack.ClientID,
			cellnet.MessageToString(userMsg))
	} else {

		// 网关没有相关的消息, 只能打出消息号
		ulog.Debugf("#%s(%s) len: %d msgid: %d cid:%d",
			dir,
			agentNodeID,
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
		))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
