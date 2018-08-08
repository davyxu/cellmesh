package service

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/cellnet/rpc"
	"reflect"
	"time"
)

type rpcRequest struct {
	conn cellnet.TCPConnector

	readyChan chan string
}

func (self *rpcRequest) Session() cellnet.Session {
	return self.conn.Session()
}

func (self *rpcRequest) Start() {
	self.conn.Start()
}

func (self *rpcRequest) IsReady() bool {
	return self.conn.(interface {
		IsReady() bool
	}).IsReady()
}

func (self *rpcRequest) Stop() {
	self.conn.Stop()
}

func (self *rpcRequest) Request(req interface{}, ackType reflect.Type, callback func(interface{})) error {

	ack, err := rpc.CallSync(self.conn, req, time.Hour)
	if err != nil {
		return err
	}

	callback(ack)

	return nil
}

func (self *rpcRequest) onMessage(ev cellnet.Event) {

	switch ev.Message().(type) {
	case *cellnet.SessionConnected: // 已经连接上
	case *cellnet.SessionClosed:
		self.readyChan <- "closed"
	}
}

func NewRPCRequestor(addr string, readyChan chan string) Requestor {

	p := peer.NewGenericPeer("tcp.SyncConnector", addr, addr, nil)

	self := &rpcRequest{
		conn:      p.(cellnet.TCPConnector),
		readyChan: readyChan,
	}

	proc.BindProcessorHandler(p, "tcp.ltv", self.onMessage)

	return self
}

/*
func (self *rpcRequest) Request(req interface{}, ackType reflect.Type, callback func(interface{})) error {

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

}
func (self *rpcRequest) onMessage(ev cellnet.Event) {

	incomingMsgType := reflect.TypeOf(ev.Message()).Elem()

	if rawFeedback, ok := self.ctxList.Load(incomingMsgType); ok {
		feedBack := rawFeedback.(chan interface{})
		feedBack <- ev.Message()
	}

	switch ev.Message().(type) {
	case *cellnet.SessionConnected: // 已经连接上
		self.readyChan <- self
	case *cellnet.SessionClosed:
		service.RemoveConnection(ev.Session().Peer().(cellnet.PeerProperty).Address())
	}
}
*/
