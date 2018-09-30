package service

import (
	"github.com/davyxu/cellnet"
	"sync"
)

type connSet struct {
	data  map[string]cellnet.Peer
	mutex sync.RWMutex
}

func (self *connSet) Add(svcid string, peer cellnet.Peer) {
	self.mutex.Lock()
	self.data[svcid] = peer
	self.mutex.Unlock()
}

func (self *connSet) Remove(svcid string) {
	self.mutex.Lock()
	delete(self.data, svcid)
	self.mutex.Unlock()
}

func (self *connSet) Get(svcid string) cellnet.Peer {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return self.data[svcid]
}

func (self *connSet) Count() int {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return len(self.data)
}

func newConnSet() *connSet {
	return &connSet{
		data: make(map[string]cellnet.Peer),
	}
}
