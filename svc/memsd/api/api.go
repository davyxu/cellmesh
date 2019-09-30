package memsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"sync"
)

type memDiscovery struct {
	config *Config

	ses      cellnet.Session
	sesGuard sync.RWMutex

	kvCache      map[string][]byte
	kvCacheGuard sync.RWMutex

	svcCache      map[string][]*discovery.ServiceDesc
	svcCacheGuard sync.RWMutex

	notifyFunc discovery.NotifyFunc

	initWg *sync.WaitGroup

	token string

	q cellnet.EventQueue
}

func (self *memDiscovery) triggerNotify(evType string, args ...interface{}) {

	if self.notifyFunc != nil {
		self.notifyFunc(evType, args...)
	}
}
func (self *memDiscovery) SetNotify(callback discovery.NotifyFunc) {

	self.notifyFunc = callback
}

func (self *memDiscovery) Start(config interface{}) {

	if config == nil {
		config = DefaultConfig()
	}

	self.config = config.(*Config)

	self.initWg = new(sync.WaitGroup)
	self.initWg.Add(1)

	self.connect(self.config.Address)

	// 等待拉取初始值
	self.initWg.Wait()
	self.initWg = nil
}

func NewDiscovery() discovery.Discovery {

	self := &memDiscovery{
		kvCache:  make(map[string][]byte),
		svcCache: make(map[string][]*discovery.ServiceDesc),
		q:        cellnet.NewEventQueue(),
	}
	self.q.EnableCapturePanic(true)
	self.q.StartLoop()

	return self
}
