package memsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/memsd/model"
	"github.com/davyxu/cellnet"
	"sync"
)

type notifyContext struct {
	stack string
	mode  string
}

type memDiscovery struct {
	config *Config

	ses      cellnet.Session
	sesGuard sync.RWMutex

	kvCache      map[string][]byte
	kvCacheGuard sync.RWMutex

	svcCache      map[string][]*discovery.ServiceDesc
	svcCacheGuard sync.RWMutex

	notifyMap sync.Map // key=mode+c value=string

	initWg *sync.WaitGroup

	token string
}

func NewDiscovery(config interface{}) discovery.Discovery {

	if config == nil {
		config = DefaultConfig()
	}

	self := &memDiscovery{
		config:   config.(*Config),
		kvCache:  make(map[string][]byte),
		svcCache: make(map[string][]*discovery.ServiceDesc),
	}

	model.Queue = cellnet.NewEventQueue()
	model.Queue.EnableCapturePanic(true)
	model.Queue.StartLoop()

	self.initWg = new(sync.WaitGroup)
	self.initWg.Add(1)

	self.connect(self.config.Address)

	// 等待拉取初始值
	self.initWg.Wait()
	self.initWg = nil

	return self
}
