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
)

type BackendMsgHooker struct {
}

// 后端服务器接收来自网关的消息
func (BackendMsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch incomingMsg := inputEvent.Message().(type) {
	case *proto.TransmitACK:

		userMsg, _, err := codec.DecodeMessage(int(incomingMsg.MsgID), incomingMsg.MsgData)
		if err != nil {
			log.Warnf("Backend msg decode failed, %s, msgid: %d", err.Error(), incomingMsg.MsgID)
			return nil
		}

		ev := &RecvMsgEvent{
			Ses:      inputEvent.Session(),
			Msg:      userMsg,
			ClientID: incomingMsg.ClientID,
		}

		outputEvent = ev

	default:
		outputEvent = inputEvent
	}

	return
}

// 后端服务器发送到网关的消息
func (BackendMsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch outgoingMsg := inputEvent.Message().(type) {
	case *proto.TransmitACK:

		if log.IsDebugEnabled() {

			writeAgentLog(inputEvent.Session(), "send", outgoingMsg)
		}
	}

	return inputEvent
}

type broadcasterHooker struct {
}

// 来自后台服务器的消息
func (broadcasterHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch incomingMsg := inputEvent.Message().(type) {
	case *proto.TransmitACK:

		rawPkt := &cellnet.RawPacket{
			MsgData: incomingMsg.MsgData,
			MsgID:   int(incomingMsg.MsgID),
		}

		if log.IsDebugEnabled() {

			writeAgentLog(inputEvent.Session(), "recv", incomingMsg)
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
	}

	return inputEvent
}

// 发送给后台服务器
func (broadcasterHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch outgoingMsg := inputEvent.Message().(type) {
	case *proto.TransmitACK:

		if log.IsDebugEnabled() {

			writeAgentLog(inputEvent.Session(), "send", outgoingMsg)
		}
	}

	return inputEvent
}

func writeAgentLog(ses cellnet.Session, dir string, ack *proto.TransmitACK) {

	if !msglog.IsMsgLogValid(int(ack.MsgID)) {
		return
	}

	peerInfo := ses.Peer().(cellnet.PeerProperty)

	userMsg, _, err := codec.DecodeMessage(int(ack.MsgID), ack.MsgData)
	if err == nil {
		log.Debugf("#agent.%s(%s)@%d len: %d %s <%d>| %s",
			dir,
			peerInfo.Name(),
			ses.ID(),
			cellnet.MessageSize(userMsg),
			cellnet.MessageToName(userMsg),
			ack.ClientID,
			cellnet.MessageToString(userMsg))
	} else {

		// 网关没有相关的消息, 只能打出消息号
		log.Debugf("#agent.%s(%s)@%d len: %d msgid: %d <%d>",
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

	// 避免默认消息日志显示本条消息
	msglog.SetMsgLogRule("proto.TransmitACK", msglog.MsgLogRule_BlackList)

	// 适用于后端服务的处理器
	proc.RegisterProcessor("svc.backend", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(
			new(link.SvcEventHooker), // 服务互联处理
			new(BackendMsgHooker),    // 网关消息处理
			new(tcp.MsgHooker)))      // tcp基础消息处理
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})

	// 适用于
	proc.RegisterProcessor("agent.backend", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(
			new(link.SvcEventHooker), // 服务互联处理
			new(broadcasterHooker),   // 网关消息处理
			new(tcp.MsgHooker)))      // tcp基础消息处理
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
