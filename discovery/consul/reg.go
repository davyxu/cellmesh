package consulsd

import (
	"github.com/davyxu/cellmicro/discovery"
	"github.com/hashicorp/consul/api"
	"time"
)

const ServiceTTL = time.Second * 5

func (self *consulDiscovery) Register(svc *discovery.ServiceDesc) error {

	log.Debugf("Register service, %s", svc.String())

	var checker api.AgentServiceCheck
	checker.CheckID = svc.ID
	checker.TTL = ServiceTTL.String()
	checker.DeregisterCriticalServiceAfter = (time.Minute + ServiceTTL).String()

	err := self.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      svc.ID,
		Name:    svc.Name,
		Address: svc.Address,
		Port:    svc.Port,
		//Check:   &checker,
	})

	if err != nil {
		return err
	}

	//go func() {
	//
	//	for {
	//
	//		self.client.Agent().UpdateTTL(svc.ID, "svc ready", "pass")
	//
	//		time.Sleep(ServiceTTL)
	//	}
	//
	//}()

	return nil
}

func (self *consulDiscovery) Deregister(svcid string) error {

	log.Debugf("Deregister service, id: %s", svcid)

	return self.client.Agent().ServiceDeregister(svcid)
}
