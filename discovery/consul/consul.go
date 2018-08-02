package consulsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
	"sync"
)

type consulDiscovery struct {
	client *api.Client

	config *api.Config

	cache      sync.Map // map[string][]*discovery.ServiceDesc
	cacheGuard sync.Mutex

	nameWatcher sync.Map //map[string]*watch.Plan
}

func consulSvcToService(s *api.ServiceEntry) *discovery.ServiceDesc {

	return &discovery.ServiceDesc{
		Name:    s.Service.Service,
		ID:      s.Service.ID,
		Address: s.Service.Address,
		Port:    s.Service.Port,
	}
}

// from github.com/micro/go-micro/registry/consul_registry.go
func (self *consulDiscovery) Query(name string) (ret []*discovery.ServiceDesc, err error) {

	log.Debugf("Query service, name: %s", name)

	result, _, err := self.client.Health().Service(name, "", false, nil)

	if err != nil {
		return nil, err
	}

	for _, s := range result {

		if s.Service.Service != name {
			continue
		}

		sd := consulSvcToService(s)

		log.Debugf("  got servcie, %s", sd.String())

		ret = append(ret, sd)
	}

	return

}

func newConsulDiscovery() discovery.Discovery {

	self := &consulDiscovery{
		config: api.DefaultConfig(),
	}

	var err error
	self.client, err = api.NewClient(self.config)

	if err != nil {
		panic(err)
	}

	self.startWatch()

	return self
}

func init() {
	discovery.Default = newConsulDiscovery()
}
