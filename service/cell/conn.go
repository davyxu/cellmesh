package cell

import (
	"errors"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"reflect"
	"sync"
	"time"
)

type cellRequest struct {
	ctxList sync.Map

	conn cellnet.TCPConnector

	readyChan chan service.Requestor
}

func (self *cellRequest) Session() cellnet.Session {
	return self.conn.Session()
}

func (self *cellRequest) Request(req interface{}, ackType reflect.Type, callback func(interface{})) error {

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

func (self *cellRequest) onMessage(ev cellnet.Event) {

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

func init() {

	service.NewRequestor = func(addr string, readyChan chan service.Requestor) service.Requestor {

		p := peer.NewGenericPeer("tcp.Connector", addr, addr, nil)

		self := &cellRequest{
			conn:      p.(cellnet.TCPConnector),
			readyChan: readyChan,
		}

		proc.BindProcessorHandler(p, "tcp.ltv", self.onMessage)

		p.Start()

		return self
	}

}
