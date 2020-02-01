package rbase

import (
	"fmt"
	robotutil "github.com/davyxu/cellmesh/svc/robot/util"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"reflect"
	"sync"
	"time"
)

type msgContext struct {
	msg      interface{}
	recvTime time.Time
	done     bool
}

func (self *msgContext) IsTimeout() bool {
	return time.Now().Sub(self.recvTime) > time.Second*30
}

type Messenger struct {
	msgList sync.Map

	waitMsg      *cellnet.MessageMeta
	waitMsgGuard sync.RWMutex

	bgProc func(interface{}) bool

	// Socket
	peerByName  map[string]cellnet.Peer
	recvTimeout time.Duration
}

func (self *Messenger) SetRecvTimeoutSec(sec int) {
	self.recvTimeout = time.Duration(sec) * time.Second
}

func (self *Messenger) VisitWaitMsg(callback func(meta *cellnet.MessageMeta) bool) {

	self.waitMsgGuard.RLock()
	callback(self.waitMsg)
	self.waitMsgGuard.RUnlock()

}

func (self *Messenger) Send(socket string, reqMsg interface{}) {

	p := self.GetPeer(socket)
	if p == nil {
		return
	}

	p.(interface {
		Session() cellnet.Session
	}).Session().Send(reqMsg)
}

func (self *Messenger) IsPeerReady(socket string) bool {

	p := self.GetPeer(socket)
	if p == nil {
		return false
	}

	return p.(cellnet.PeerReadyChecker).IsReady()
}

func (self *Messenger) GetPeer(peerName string) cellnet.Peer {

	p := self.peerByName[peerName]
	if p == nil {
		panic("invalid peer")
	}

	return p
}

func (self *Messenger) SetPeer(peerName string, peer cellnet.Peer) {
	self.peerByName[peerName] = peer
}

func (self *Messenger) SetBackgroundRecv(callback func(msg interface{}) bool) {

	self.bgProc = callback
}

func (self *Messenger) backgroundProc() {
	var copyList []*cellnet.MessageMeta
	self.msgList.Range(func(key, value interface{}) bool {

		meta := key.(*cellnet.MessageMeta)
		ctx := value.(*msgContext)

		var needDelete bool
		if ctx.done || ctx.IsTimeout() { // 处理 和超时
			needDelete = true
		} else if self.bgProc(ctx.msg) { // 后台处理
			needDelete = true
		}

		if needDelete {
			copyList = append(copyList, meta)
		}

		return true
	})

	for _, meta := range copyList {
		self.msgList.Delete(meta)
	}
}

func (self *Messenger) recvMulti(expectMeta *cellnet.MessageMeta, callback func(interface{})) {

	var copyList []interface{}
	self.msgList.Range(func(key, value interface{}) bool {

		meta := key.(*cellnet.MessageMeta)
		ctx := value.(*msgContext)

		if meta == expectMeta {
			copyList = append(copyList, ctx.msg)
			ctx.done = true
		}

		return true
	})

	for _, msg := range copyList {
		callback(msg)
	}

	self.backgroundProc()

}

func (self *Messenger) Recv(msgName string) (ret interface{}) {

	expectMeta := cellnet.MessageMetaByFullName(msgName)

	self.waitMsgGuard.Lock()
	self.waitMsg = expectMeta
	self.waitMsgGuard.Unlock()

	done := false

	beginTime := time.Now()
	for !done {

		self.recvMulti(expectMeta, func(msg interface{}) {

			if !done {
				ret = msg
				done = true
			}
		})

		time.Sleep(time.Millisecond * 100)

		if time.Now().After(beginTime.Add(self.recvTimeout)) {
			panic(fmt.Errorf("RecvTimeout, msgName: %s  stack: %s", msgName, util.StackToString(3)))
		}
	}

	self.waitMsgGuard.Lock()
	self.waitMsg = nil
	self.waitMsgGuard.Unlock()

	return
}

var FastExec bool

func (self *Messenger) Sleep() {

	if FastExec {
		return
	}

	// 模拟延迟
	delay := robotutil.RandRangeInt32(100, 1500)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func (self *Messenger) SleepRange(beginMS, endMS int32) {

	if FastExec {
		return
	}

	// 模拟延迟
	delay := robotutil.RandRangeInt32(beginMS, endMS)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func (self *Messenger) RecvAsync(callback interface{}) {

	funcType := reflect.TypeOf(callback)
	replyType := funcType.In(0)

	expectMeta := cellnet.MessageMetaByType(replyType)

	self.recvMulti(expectMeta, func(msg interface{}) {
		reflect.ValueOf(callback).Call([]reflect.Value{reflect.ValueOf(msg)})
	})
}

func (self *Messenger) OnEvent(inputEvent cellnet.Event) {

	meta := cellnet.MessageMetaByMsg(inputEvent.Message())

	if meta != nil {
		self.msgList.Store(meta, &msgContext{
			recvTime: time.Now(),
			msg:      inputEvent.Message(),
		})
	}
}

func (self *Messenger) Init() {
	self.peerByName = make(map[string]cellnet.Peer)
}
