package rpc

import (
	"fmt"
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/link"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/ulog"
)

type MsgHooker struct {
}

func decodePayload(rpcMsg *proto.HubTransmitACK) (userMsg, userPt interface{}, err error) {
	userMsg, _, err = codec.DecodeMessage(int(rpcMsg.MsgID), rpcMsg.MsgData)
	if err != nil {
		err = fmt.Errorf("rpc decode failed, %s", err)
		return
	}

	if rpcMsg.PassThroughType != "" {
		userPt, err = decodePassthrough(rpcMsg.PassThroughData, rpcMsg.PassThroughType)
		if err != nil {
			err = fmt.Errorf("rpc decode failed, %s", err)
			return
		}
	}

	return
}

func (self MsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	rpcMsg, ok := inputEvent.Message().(*proto.HubTransmitACK)
	if !ok {
		return inputEvent
	}

	switch rpcMsg.Mode {
	case TransmitMode_Reply: // 客户端接收并回调

		if rpcMsg.TgtSvcID != fx.LocalSvcID {
			ulog.Warnf("recv invalid rpc target svcid: %s msgid: %d", rpcMsg.TgtSvcID, rpcMsg.MsgID)
			return inputEvent
		}

		userMsg, userPt, err := decodePayload(rpcMsg)
		if err != nil {
			ulog.Errorf("%s", err)
			return inputEvent
		}

		if ulog.IsLevelEnabled(ulog.DebugLevel) {
			peerInfo := inputEvent.Session().Peer().(cellnet.PeerProperty)

			ulog.Debugf("#rpc.recv(%s)@%d len: %d %s | %s",
				peerInfo.Name(),
				inputEvent.Session().ID(),
				cellnet.MessageSize(userMsg),
				cellnet.MessageToName(userMsg),
				cellnet.MessageToString(userMsg))
		}

		mgr := fetchManager(inputEvent.Session())
		req := mgr.Get(rpcMsg.CallID)

		if req != nil {

			req.onRespond(userMsg, userPt, err)
		} else {
			ulog.Warnf("rpc respond not hit, id: %d, msgid: %d srcSvcID: %s", rpcMsg.CallID, rpcMsg.MsgID, rpcMsg.SrcSvcID)
		}

	case TransmitMode_RequestNotify: // 服务器接收并回应
		// 接收到本进程应该接收的消息
		if rpcMsg.TgtSvcID == fx.LocalSvcID {

			userMsg, userPt, err := decodePayload(rpcMsg)
			if err != nil {
				ulog.Errorf("%s", err)
				return inputEvent
			}

			if ulog.IsLevelEnabled(ulog.DebugLevel) {
				peerInfo := inputEvent.Session().Peer().(cellnet.PeerProperty)

				ulog.Debugf("#rpc.recv(%s)@%d len: %d %s | %s",
					peerInfo.Name(),
					inputEvent.Session().ID(),
					cellnet.MessageSize(userMsg),
					cellnet.MessageToName(userMsg),
					cellnet.MessageToString(userMsg))
			}

			return &RecvMsgEvent{
				ses:      inputEvent.Session(),
				Msg:      userMsg,
				callid:   rpcMsg.CallID,
				srcSvcID: rpcMsg.SrcSvcID,
				recvPt:   userPt,
			}

		} else {
			// 路由到远程服务器
			remoteSvcID := link.GetLink(rpcMsg.TgtSvcID)
			if remoteSvcID == nil {
				ulog.Warnf("rpc target not found, svcid: '%s', msgid: %d", rpcMsg.TgtSvcID, rpcMsg.MsgID)
				return inputEvent
			}

			remoteSvcID.Send(rpcMsg)

			if ulog.IsLevelEnabled(ulog.DebugLevel) {

				userMsg, _, err := decodePayload(rpcMsg)
				if err != nil {
					ulog.Errorf("%s", err)
					return inputEvent
				}

				peerInfo := inputEvent.Session().Peer().(cellnet.PeerProperty)

				ulog.Debugf("#rpc.relay(%s)@%d len: %d %s | %s",
					peerInfo.Name(),
					inputEvent.Session().ID(),
					cellnet.MessageSize(userMsg),
					cellnet.MessageToName(userMsg),
					cellnet.MessageToString(userMsg))
			}
		}
	}

	return inputEvent
}

func (self MsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	rpcMsg, ok := inputEvent.Message().(*proto.HubTransmitACK)
	if !ok {
		return inputEvent
	}

	if ulog.IsLevelEnabled(ulog.DebugLevel) {
		userMsg, _, err := decodePayload(rpcMsg)
		if err != nil {
			ulog.Errorf("%s", err)
			return inputEvent
		}

		peerInfo := inputEvent.Session().Peer().(cellnet.PeerProperty)

		ulog.Debugf("#rpc.send(%s)@%d len: %d %s | %s",
			peerInfo.Name(),
			inputEvent.Session().ID(),
			cellnet.MessageSize(userMsg),
			cellnet.MessageToName(userMsg),
			cellnet.MessageToString(userMsg))
	}

	return inputEvent
}
