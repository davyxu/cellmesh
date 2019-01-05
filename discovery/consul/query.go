package consulsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"time"
)

func (self *consulDiscovery) updateFromConsul() {

	newCache := map[string][]*discovery.ServiceDesc{}
	for _, desc := range self.QueryAll() {

		list := newCache[desc.Name]
		list = append(list, desc)
		newCache[desc.Name] = list
	}

	self.svcCacheFullGuard.Lock()
	self.svcCacheFull = newCache
	self.svcCacheFullGuard.Unlock()
}

func (self *consulDiscovery) queryFromCache(name string) (ret []*discovery.ServiceDesc) {

	self.svcCacheFullGuard.RLock()
	list := self.svcCacheFull[name]

	ret = make([]*discovery.ServiceDesc, len(list))
	copy(ret, list)
	self.svcCacheFullGuard.RUnlock()

	return
}

func (self *consulDiscovery) Query(name string) (ret []*discovery.ServiceDesc) {

	//if raw, ok := self.svcCache.Load(name); ok {
	//	ret = raw.([]*discovery.ServiceDesc)
	//}

	self.cacheDirtyGuard.Lock()
	if self.cacheDirty {
		self.updateFromConsul()
		self.cacheDirty = true
	}
	self.cacheDirtyGuard.Unlock()

	return self.queryFromCache(name)
}

// from github.com/micro/go-micro/registry/consul_registry.go
func (self *consulDiscovery) directQuery(name string) (ret []*discovery.ServiceDesc, err error) {

	result, _, err := self.client.Health().Service(name, "", false, nil)

	if err != nil {
		return nil, err
	}

	for _, s := range result {

		if s.Service.Service != name {
			continue
		}

		if isServiceHealth(s) {

			sd := consulSvcToService(s.Service)

			log.Debugf("  got servcie, %s", sd.String())

			ret = append(ret, sd)
		}

	}

	return

}

func (self *consulDiscovery) QueryAll() (ret []*discovery.ServiceDesc) {

	svc, err := self.client.Agent().Services()
	if err != nil {
		return
	}

	for _, detail := range svc {
		ret = append(ret, consulSvcToService(detail))
	}

	return
}

func (self *consulDiscovery) RegisterNotify(mode string) (ret chan struct{}) {

	ret = make(chan struct{}, 10)

	switch mode {
	case "add":
		self.notifyMap.Store(ret, struct{}{})
	case "remove":
	}

	return
}

func (self *consulDiscovery) DeregisterNotify(mode string, c chan struct{}) {

	switch mode {
	case "add":
		self.notifyMap.Store(c, nil)
	case "remove":
	}
}

func (self *consulDiscovery) OnCacheUpdated(eventName string, desc *discovery.ServiceDesc) {

	self.svcCacheFullGuard.Lock()
	self.cacheDirty = true
	self.svcCacheFullGuard.Unlock()

	switch eventName {
	case "add":

		var notifyList []chan struct{}
		var removeList []chan struct{}

		// 将列表拷贝出来，避免互锁
		self.notifyMap.Range(func(key, value interface{}) bool {
			if value != nil {
				notifyList = append(notifyList, key.(chan struct{}))
			} else {
				removeList = append(removeList, key.(chan struct{}))
			}

			return true
		})

		// 通知
		for _, raw := range notifyList {
			notify(raw)
		}

		// 删除已经解注册的
		for _, raw := range removeList {
			self.notifyMap.Delete(raw)
		}

	case "remove":

	}

}

func notify(c chan struct{}) {

	select {
	case c <- struct{}{}:
	case <-time.After(10 * time.Second):
		// 接收通知阻塞太久，或者没有释放侦听的channel
		log.Errorf("addNotify timeout, not free?")
	}

}
