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

	self.notifyGuard.Lock()
	switch mode {
	case "add":
		self.addNotify = append(self.addNotify, ret)
	case "remove":
		self.removeNotify = append(self.removeNotify, ret)
	}
	self.notifyGuard.Unlock()

	return
}

func (self *consulDiscovery) DeregisterNotify(mode string, c chan struct{}) {

	self.notifyGuard.Lock()
	switch mode {
	case "add":
		for index, n := range self.addNotify {
			if n == c {
				self.addNotify = append(self.addNotify[:index], self.addNotify[index+1:]...)
				break
			}
		}
	case "remove":
		for index, n := range self.removeNotify {
			if n == c {
				self.removeNotify = append(self.removeNotify[:index], self.removeNotify[index+1:]...)
				break
			}
		}
	}
	self.notifyGuard.Unlock()

}

func (self *consulDiscovery) OnCacheUpdated(eventName string, desc *discovery.ServiceDesc) {

	self.notifyGuard.RLock()
	switch eventName {
	case "add":
		//log.Debugf("Add service '%s'", desc.ID)
		notify(self.addNotify)

	case "remove":
		//log.Debugf("Remove service '%s'", desc.ID)
		notify(self.removeNotify)
	}

	self.notifyGuard.RUnlock()
}

func notify(clist []chan struct{}) {
	for _, n := range clist {

		select {
		case n <- struct{}{}:
		case <-time.After(5 * time.Second):
			log.Errorf("addNotify timeout, not free?")

		}

	}
}
