package consulsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
	"time"
)

// 心跳时间
const ServiceTTL = 5 * time.Second

// 心跳超时
const ServiceTTLTimeout = 7 * time.Second

// Consul移除服务时间，Consul要求必须在1分钟后才能删除TTL超时的服务
const RemoveServiceTimeout = time.Minute + 5*time.Second

func (self *consulDiscovery) Register(svc *discovery.ServiceDesc) error {

	log.Debugf("Register service, %s", svc.String())

	var checker api.AgentServiceCheck
	checker.CheckID = svc.ID
	checker.TTL = ServiceTTLTimeout.String()
	checker.DeregisterCriticalServiceAfter = RemoveServiceTimeout.String()

	err := self.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      svc.ID,
		Name:    svc.Name,
		Address: svc.Host,
		Port:    svc.Port,
		Tags:    svc.Tags,
		Check:   &checker,
		Meta:    svc.Meta,
	})

	if err != nil {
		return err
	}

	localSvc := newLocalService(svc, self.client.Agent())

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
	}

	log.Debugf("Deregister service, id: %s", svcid)

	return self.client.Agent().ServiceDeregister(svcid)
}
