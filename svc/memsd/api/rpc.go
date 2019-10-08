package memsd

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/rpc"
)

func (self *memDiscovery) RPCSession() (ses cellnet.Session) {
	self.sesGuard.RLock()
	ses = self.ses
	self.sesGuard.RUnlock()
	return
}

// callback =func(ack *YouMsgACK)
func (self *memDiscovery) remoteCall(req interface{}, callback interface{}) {
	rpc.CallSyncType(self.RPCSession(), req, self.config.RequestTimeout, callback)
}
