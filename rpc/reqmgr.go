package rpc

import (
	"github.com/davyxu/cellnet"
	"sync"
	"sync/atomic"
)

type RequestManager struct {
	reqByID sync.Map
	idAcc   int64
}

func (self *RequestManager) Get(id int64) *Request {
	if v, ok := self.reqByID.Load(id); ok {
		return v.(*Request)
	}

	return nil
}

func (self *RequestManager) Add(req *Request) {
	req.id = atomic.AddInt64(&self.idAcc, 1)
	self.reqByID.Store(req.id, req)
}

func (self *RequestManager) Remove(req *Request) {
	self.reqByID.Delete(req.id)
}

func NewRequestManager() *RequestManager {
	return &RequestManager{}
}

const (
	peerContextKey_reqmgr = "rpcmgr"
)

func fetchManager(ses cellnet.Session) *RequestManager {
	if ctx, ok := ses.Peer().(cellnet.ContextSet); ok {
		if v, ok := ctx.GetContext(peerContextKey_reqmgr); ok {
			return v.(*RequestManager)
		}

		mgr := NewRequestManager()
		ctx.SetContext(peerContextKey_reqmgr, mgr)
		return mgr
	} else {
		panic("peer not support context set")
	}
}
