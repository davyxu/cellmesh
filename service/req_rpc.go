package service

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/cellnet/rpc"
	"reflect"
	"sync"
	"time"
)

type rpcRequest struct {
	conn cellnet.TCPConnector

	stop sync.WaitGroup
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
		self.stop.Done()
	}
}

func (self *rpcRequest) WaitStop() {

	self.stop.Wait()
}

func NewRPCRequestor(addr string) Requestor {

	p := peer.NewGenericPeer("tcp.SyncConnector", addr, addr, nil)

	self := &rpcRequest{
		conn: p.(cellnet.TCPConnector),
	}

	self.stop.Add(1)

	proc.BindProcessorHandler(p, "tcp.ltv", self.onMessage)

	return self
}
