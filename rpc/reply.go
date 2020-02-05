package rpc

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/ulog"
)

const (
	TransmitMode_RequestNotify = 0
	TransmitMode_Reply         = 1
)

type RecvMsgEvent struct {
	ses      cellnet.Session
	Msg      interface{}
	callid   int64
	srcSvcID string
	recvPt   interface{}
	replyPt  interface{}
}

func (self *RecvMsgEvent) Session() cellnet.Session {
	return self.ses
}

func (self *RecvMsgEvent) Message() interface{} {
	return self.Msg
}

func (self *RecvMsgEvent) Queue() cellnet.EventQueue {
	return self.ses.Peer().(interface {
		Queue() cellnet.EventQueue
	}).Queue()
}

func (self *RecvMsgEvent) WithPassThrough(data interface{}) {
	self.replyPt = data
}

func (self *RecvMsgEvent) Reply(msg interface{}) {

	var (
		ack proto.HubTransmitACK
		err error
	)

	ack.MsgData, ack.MsgID, err = saveMessage(msg)
	if err != nil {
		ulog.Errorf("rpc reply message encode error: %s", err)
		return
	}

	ack.Mode = TransmitMode_Reply
	ack.CallID = self.callid
	ack.SrcSvcID = fx.LocalSvcID
	ack.TgtSvcID = self.srcSvcID

	if self.replyPt != nil {
		ack.PassThroughData, ack.PassThroughType, err = savePassthrough(self.replyPt)
	}

	self.ses.Send(&ack)
}
