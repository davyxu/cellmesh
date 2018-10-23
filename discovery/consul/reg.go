package consulsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
)

func (self *consulDiscovery) Register(svc *discovery.ServiceDesc) error {

	//log.Debugf("Register service, %s", svc.String())

	var checker api.AgentServiceCheck
	checker.CheckID = svc.ID
	checker.TTL = self.config.ServiceTTLTimeOut.String()
	checker.DeregisterCriticalServiceAfter = self.config.RemoveServiceTimeout.String()

	err := self.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:                svc.ID,
		Name:              svc.Name,
		Address:           svc.Host,
		Port:              svc.Port,
		Tags:              svc.Tags,
		Check:             &checker,
		Meta:              svc.Meta,
		EnableTagOverride: true,
	})

	if err != nil {
		return err
	}

	localSvc := newLocalService(self, svc, self.client.Agent())

	self.localSvc.Store(svc.ID, localSvc)

	return nil
}

// 重新注册
func (self *consulDiscovery) Recover() {

	self.localSvc.Range(func(key, value interface{}) bool {
		svc := value.(*localService)

		return self.Register(svc.Desc) == nil
	})

}

func (self *consulDiscovery) Deregister(svcid string) error {

	if v, ok := self.localSvc.Load(svcid); ok {
		localsvc := v.(*localService)

		localsvc.Stop()

		self.localSvc.Delete(svcid)
	}

	//log.Debugf("Deregister service, id: %s", svcid)

	return self.client.Agent().ServiceDeregister(svcid)
}

func (self *consulDiscovery) SetChecker(svcid string, checkerFunc discovery.CheckerFunc) {

	if v, ok := self.localSvc.Load(svcid); ok {
		localsvc := v.(*localService)

		localsvc.checkerFunc = checkerFunc
	}
}
