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

		mgr := fetchManager(inputEvent.Session())
		req := mgr.Get(rpcMsg.CallID)

		var (
			err     error
			recvMsg interface{}
			recvPt  interface{}
		)

		recvMsg, _, err = codec.DecodeMessage(int(rpcMsg.MsgID), rpcMsg.MsgData)
		if err != nil {
			req.onRespond(recvMsg, recvPt, fmt.Errorf("rpc decode failed, %s", err))
			return inputEvent
		}

		if rpcMsg.PassThroughType != "" {
			recvPt, err = loadPassthrough(rpcMsg.PassThroughData, rpcMsg.PassThroughType)
			if err != nil {
				req.onRespond(recvMsg, recvPt, fmt.Errorf("rpc decode failed, %s", err))
				return inputEvent
			}
		}
	case TransmitMode_RequestNotify: // 服务器接收并回应
		// 接收到本进程应该接收的消息
		if rpcMsg.TgtSvcID == fx.LocalSvcID {

			var (
				err     error
				recvMsg interface{}
				recvPt  interface{}
			)
			recvMsg, _, err = codec.DecodeMessage(int(rpcMsg.MsgID), rpcMsg.MsgData)
			if err != nil {
				ulog.Errorf("rpc decode failed, %s", err)
				return inputEvent
			}

			if rpcMsg.PassThroughType != "" {
				recvPt, err = loadPassthrough(rpcMsg.PassThroughData, rpcMsg.PassThroughType)
				if err != nil {
					ulog.Errorf("rpc decode failed, %s", err)
					return inputEvent
				}
			}
			return &RecvMsgEvent{
				ses:      inputEvent.Session(),
				Msg:      recvMsg,
				callid:   rpcMsg.CallID,
				srcSvcID: rpcMsg.SrcSvcID,
				recvPt:   recvPt,
			}

		} else {
			// 路由到远程服务器
			remoteSvcID := link.GetLink(rpcMsg.TgtSvcID)
			if remoteSvcID == nil {
				ulog.Warnf("rpc target not found, svcid: '%s', msgid: %d", rpcMsg.TgtSvcID, rpcMsg.MsgID)
				return inputEvent
			}

			remoteSvcID.Send(rpcMsg)
		}
	}

	return inputEvent
}

func (self MsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}
