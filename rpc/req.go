package rpc

import (
	"errors"
	"fmt"
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/link"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"time"
)

var ErrTimeout = errors.New("rpc time out")

type Respond struct {
	Passsthrough interface{}
	Message      interface{}
	Error        error
}

type Request struct {
	id      int64
	pt      interface{}
	timeout time.Duration
	ses     cellnet.Session
	err     error

	ch       chan *Respond
	callback func(*Respond)
}

func (self *Request) Request(msg interface{}) *Request {

	var (
		ack proto.HubTransmitACK
		err error
	)

	ack.MsgData, ack.MsgID, err = saveMessage(msg)
	if err != nil {
		self.err = fmt.Errorf("rpc request message encode error: %w", err)
		return self
	}

	if self.pt != nil {
		ack.PassThroughData, ack.PassThroughType, err = savePassthrough(self.pt)
	}

	ack.Mode = TransmitMode_RequestNotify
	ack.SrcSvcID = fx.LocalSvcID
	ack.TgtSvcID = link.GetLinkSvcID(self.ses)

	fetchManager(self.ses).Add(self)

	self.ses.Send(&ack)

	return self
}

func (self *Request) WithTimeout(dur time.Duration) *Request {
	self.timeout = dur
	return self
}

// 透传消息
func (self *Request) WithPassThrough(msg interface{}) *Request {
	self.pt = msg
	return self
}

func (self *Request) Recv(callback func(*Respond)) {
	self.callback = callback

	if self.err != nil {
		cellnet.SessionQueuedCall(self.ses, func() {
			callback(&Respond{
				Error: self.err,
			})
		})
	}
}

func (self *Request) RecvWait(callback func(*Respond)) {

	if self.err != nil {
		callback(&Respond{
			Error: self.err,
		})
	} else {
		self.ch = make(chan *Respond)

		select {
		case r := <-self.ch:
			callback(r)
		case <-time.After(self.timeout):
			callback(&Respond{
				Error: ErrTimeout,
			})
		}
	}
}

func (self *Request) onRecv(msg, pt interface{}, err error) {

	resp := &Respond{
		Passsthrough: pt,
		Message:      msg,
		Error:        err,
	}

	if self.ch != nil {
		self.ch <- resp
	} else if self.callback != nil {
		cellnet.SessionQueuedCall(self.ses, func() {
			self.callback(resp)
		})
	}
}

func NewRequest(ses cellnet.Session) *Request {

	self := &Request{
		ses:     ses,
		timeout: time.Second * 3,
	}

	return self
}

func saveMessage(msg interface{}) (data []byte, msgID uint32, err error) {
	var (
		msgMeta *cellnet.MessageMeta
	)

	data, msgMeta, err = codec.EncodeMessage(msg, nil)
	if err != nil {
		return
	}
	msgID = uint32(msgMeta.ID)

	return
}
