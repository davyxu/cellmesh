package consulsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
	"time"
)

const ServiceTTL = time.Second * 5

func (self *consulDiscovery) Register(svc *discovery.ServiceDesc) error {

	log.Debugf("Register service, %s", svc.String())

	var checker api.AgentServiceCheck
	checker.CheckID = svc.ID
	checker.TTL = ServiceTTL.String()

	// Consul要求必须在1分钟后才能删除TTL超时的服务
	checker.DeregisterCriticalServiceAfter = (time.Minute + ServiceTTL).String()

	err := self.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      svc.ID,
		Name:    svc.Name,
		Address: svc.Address,
		Port:    svc.Port,
		Check:   &checker,
	})

	if err != nil {
		return err
	}

	localSvc := newLocalService(svc.ID, self.client.Agent())

	self.localSvc.Store(svc.ID, localSvc)

	return nil
}

func (self *consulDiscovery) Deregister(svcid string) error {

	if v, ok := self.localSvc.Load(svcid); ok {
		localsvc := v.(*localService)

		localsvc.Stop()
	}

	log.Debugf("Deregister service, id: %s", svcid)

	return self.client.Agent().ServiceDeregister(svcid)
	//return nil
}
