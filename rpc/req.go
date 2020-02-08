package rpc

import (
	"errors"
	"fmt"
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/link"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"reflect"
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
	req     interface{} // 存档留作报告错误用
	pt      interface{}
	timeout time.Duration
	ses     cellnet.Session
	err     error
	mgr     *RequestManager

	ch       chan *Respond
	callback func(*Respond)
}

func (self *Request) Request(msg interface{}) *Request {

	var (
		ack proto.HubTransmitACK
		err error
	)
	self.req = msg

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

func (self *Request) Recv(callback func(resp *Respond)) {
	self.callback = callback

	if self.err != nil {
		cellnet.SessionQueuedCall(self.ses, func() {
			callback(&Respond{
				Error: self.err,
			})
		})
	} else {
		time.AfterFunc(self.timeout, func() {
			if self.mgr.Get(self.id) != nil {
				cellnet.SessionQueuedCall(self.ses, func() {
					callback(&Respond{
						Error: fmt.Errorf("request %s failed, %s", self.msgName(), ErrTimeout),
					})
				})
			}
		})

	}
}

func (self *Request) msgName() string {
	t := reflect.TypeOf(self.req)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

func (self *Request) RecvWait(callback func(resp *Respond)) {

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
				Error: fmt.Errorf("request %s failed, %s", self.msgName(), ErrTimeout),
			})
		}
	}
}

// 接收到RPC回应
func (self *Request) onRespond(msg, pt interface{}, err error) {

	resp := &Respond{
		Passsthrough: pt,
		Message:      msg,
		Error:        err,
	}

	self.mgr.Remove(self)

	if self.ch != nil {
		self.ch <- resp
	} else if self.callback != nil {
		cellnet.SessionQueuedCall(self.ses, func() {
			self.callback(resp)
		})
	}

}

func New(ses cellnet.Session) *Request {

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
