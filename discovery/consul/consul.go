package consulsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
	"sync"
	"time"
)

type consulDiscovery struct {
	client *api.Client

	config *Config

	// 与consul的服务保持实时同步
	svcCache       sync.Map // map[string][]*discovery.ServiceDesc
	svcCacheGuard  sync.Mutex
	svcNameWatcher sync.Map //map[string]*watch.Plan

	kvCache sync.Map // map[string]*KVPair

	// 本地服务的心跳
	localSvc sync.Map // map[string]*localService

	notifyMap sync.Map // key=mode+c value=chan struct{}

	// 带缓冲kv
	metaByKey sync.Map //map[string]*cacheValue
}

func (self *consulDiscovery) Raw() interface{} {
	return self.client
}

func (self *consulDiscovery) WaitReady() {

	for {

		_, _, err := self.client.Health().Service("consul", "", false, nil)

		if err == nil {
			break
		}

		log.Errorln(err)

		time.Sleep(time.Second * 2)
	}
}

func NewDiscovery(config interface{}) discovery.Discovery {

	if config == nil {
		config = DefaultConfig()
	}

	self := &consulDiscovery{
		config: config.(*Config),
	}

	var err error
	self.client, err = api.NewClient(self.config.Config)

	if err != nil {
		panic(err)
	}

	self.WaitReady()

	self.startWatchService()
	self.startWatchKV()

	return self
}
