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

		if isMeshServiceHealth(s) {

			sd := consulSvcToService(s)

			log.Debugf("  got servcie, %s", sd.String())

			ret = append(ret, sd)
		}

	}

	return

}

func (self *consulDiscovery) WaitAdded() {

	var data []interface{}
	self.pipe.Pick(&data)

	for _, d := range data {

		if d.(string) == "add" {
			return
		}
	}
}

func (self *consulDiscovery) OnCacheUpdated(eventName string, desc *discovery.ServiceDesc) {

	if self.firstUpdate != nil {
		self.firstUpdate <- struct{}{}
	}

	switch eventName {
	case "add":
		log.Debugf("Add service '%s'", desc.ID)

	case "remove":
		log.Debugf("Remove service '%s'", desc.ID)
	}

	self.pipe.Add(eventName)
}
