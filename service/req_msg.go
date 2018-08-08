package service

import (
	"errors"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"reflect"
	"sync"
	"time"
)

type msgRequest struct {
	conn cellnet.TCPConnector

	readyChan chan string

	ctxList sync.Map
}

func (self *msgRequest) Session() cellnet.Session {
	return self.conn.Session()
}

func (self *msgRequest) Start() {
	self.conn.Start()
}

func (self *msgRequest) IsReady() bool {
	return self.conn.(interface {
		IsReady() bool
	}).IsReady()
}

func (self *msgRequest) Stop() {
	self.conn.Stop()
}

func (self *msgRequest) Request(req interface{}, ackType reflect.Type, callback func(interface{})) error {

	self.conn.Session().Send(req)

	feedBack := make(chan interface{})

	self.ctxList.Store(ackType, feedBack)

	defer self.ctxList.Delete(ackType)

	select {
	case ack := <-feedBack:
		callback(ack)

		return nil
	case <-time.After(time.Second):

		return errors.New("recv time out")
	}
	return nil
}

func (self *msgRequest) onMessage(ev cellnet.Event) {

	incomingMsgType := reflect.TypeOf(ev.Message()).Elem()

	if rawFeedback, ok := self.ctxList.Load(incomingMsgType); ok {
		feedBack := rawFeedback.(chan interface{})
		feedBack <- ev.Message()
	}

	switch ev.Message().(type) {
	case *cellnet.SessionConnected: // 已经连接上
	case *cellnet.SessionClosed:
		if self.readyChan != nil {
			self.readyChan <- "closed"
		}

	}
}

func NewMsgRequestor(addr string, readyChan chan string) Requestor {

	p := peer.NewGenericPeer("tcp.SyncConnector", addr, addr, nil)

	self := &msgRequest{
		conn:      p.(cellnet.TCPConnector),
		readyChan: readyChan,
	}

	proc.BindProcessorHandler(p, "tcp.ltv", self.onMessage)

	return self
}
