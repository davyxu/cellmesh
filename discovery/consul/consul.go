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
	cache      sync.Map // map[string][]*discovery.ServiceDesc
	cacheGuard sync.Mutex

	nameWatcher sync.Map //map[string]*watch.Plan
	localSvc    sync.Map // map[string]*localService

	ready bool

	notifyGuard  sync.RWMutex
	addNotify    []chan struct{}
	removeNotify []chan struct{}

	// 带缓冲kv
	metaByKey sync.Map //map[string]*cacheValue
}

func (self *consulDiscovery) Raw() interface{} {
	return self.client
}

// 检查Consul自己挂掉
func (self *consulDiscovery) consulChecker() {

	self.ready = true

	for {

		_, _, err := self.client.Health().Service("consul", "", false, nil)

		var thisReady bool

		if err == nil {
			thisReady = true
		}

		switch {
		case self.ready == true && thisReady == false: // 宕机
			log.Warnf("Consul not reachable...")

		case self.ready == false && thisReady == true: // 恢复
			log.Warnf("Consul recovered, reregistering service...")

			// 恢复注册，虽然Consul有持久化，但是在宕机期间有注册时，还是需要重新注册
			self.Recover()
		}

		if thisReady != self.ready {
			self.ready = thisReady
		}

		time.Sleep(time.Second)

	}
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

	self.startWatch()

	return self
}
