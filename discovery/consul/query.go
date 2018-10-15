package consulsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"time"
)

func (self *consulDiscovery) startRefresh() {

	for {

		time.Sleep(time.Second * 3)

		newList := self.QueryAll()

		self.fullCacheGuard.Lock()
		self.fullCache = newList
		self.fullCacheGuard.Unlock()

		self.OnCacheUpdated("add", nil)
	}

}

func (self *consulDiscovery) queryFromCache(name string) (ret []*discovery.ServiceDesc) {
	self.fullCacheGuard.RLock()
	for _, sd := range self.fullCache {
		if sd.Name == name {
			ret = append(ret, sd)
		}
	}
	self.fullCacheGuard.RUnlock()

	return
}

func (self *consulDiscovery) Query(name string) (ret []*discovery.ServiceDesc) {

	if raw, ok := self.cache.Load(name); ok {
		ret = raw.([]*discovery.ServiceDesc)
	}
	return

	//return self.queryFromCache(name)
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

	ret = make(chan struct{})

	switch mode {
	case "add":
		self.notifyMap.Store(ret, ret)
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

	switch eventName {
	case "add":
		var needToDelete []chan struct{}
		self.notifyMap.Range(func(key, value interface{}) bool {
			c := key.(chan struct{})
			notify(c)
			needToDelete = append(needToDelete, c)

			return true
		})

		for _, c := range needToDelete {
			close(c)
			self.notifyMap.Delete(c)
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
