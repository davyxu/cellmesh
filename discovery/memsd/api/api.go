package memsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/memsd/model"
	"github.com/davyxu/cellmesh/discovery/memsd/proto"
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

	notifyMap sync.Map // key=mode+c value=string
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

	self.connect(self.config.Address)

	var wg sync.WaitGroup
	wg.Add(1)

	pullErr := self.remoteCall(&proto.PullValueREQ{}, func(ack *proto.PullValueACK) {

		model.Queue.Post(func() {

			// Pull的消息还要在queue里处理，这里确认处理完成后才算初始化完成
			wg.Done()
		})
	})

	if pullErr != nil {
		log.Errorf("Pull value failed: %s", pullErr.Error())
	} else {

		wg.Wait()
	}

	return self
}
