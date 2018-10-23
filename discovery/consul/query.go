package consulsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"time"
)

func (self *consulDiscovery) Query(name string) (ret []*discovery.ServiceDesc) {

	if raw, ok := self.svcCache.Load(name); ok {
		ret = raw.([]*discovery.ServiceDesc)
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
