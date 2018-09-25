package consulsd

import (
	"github.com/davyxu/cellmesh/discovery"
)

func (self *consulDiscovery) Query(name string) (ret []*discovery.ServiceDesc, err error) {

	log.Debugf("Query service, name: %s", name)

	if raw, ok := self.cache.Load(name); ok {
		ret = raw.([]*discovery.ServiceDesc)
	}

	return
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

			sd := consulSvcToService(s)

			log.Debugf("  got servcie, %s", sd.String())

			ret = append(ret, sd)
		}

	}

	return

}

func (self *consulDiscovery) RegisterNotify(mode string) (ret chan struct{}) {

	ret = make(chan struct{})

	switch mode {
	case "add":
		self.addNotify = append(self.addNotify, ret)
	case "remove":
		self.removeNotify = append(self.removeNotify, ret)
	}

	return
}

func (self *consulDiscovery) DeregisterNotify(mode string, c chan struct{}) {

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

}

func (self *consulDiscovery) OnCacheUpdated(eventName string, desc *discovery.ServiceDesc) {

	switch eventName {
	case "add":
		log.Debugf("Add service '%s'", desc.ID)

		for _, n := range self.addNotify {
			n <- struct{}{}
		}

	case "remove":
		log.Debugf("Remove service '%s'", desc.ID)
	}
}
