package consulsd

import (
	"context"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
	"time"
)

// 本地服务更新TTL
type localService struct {
	sd     *consulDiscovery
	Desc   *discovery.ServiceDesc
	Cancel context.CancelFunc

	ctx context.Context

	agent *api.Agent
}

func (self *localService) Update() {

	for {

		select {
		case <-self.ctx.Done():
			return
		default:

			self.agent.UpdateTTL(self.Desc.ID, self.Desc.ID, "pass")

			time.Sleep(self.sd.config.ServiceTTL)
		}
	}
}

func (self *localService) Stop() {
	self.Cancel()
}

func newLocalService(sd *consulDiscovery, svc *discovery.ServiceDesc, agent *api.Agent) *localService {

	ctx, cancel := context.WithCancel(context.Background())

	self := &localService{
		Desc:   svc,
		Cancel: cancel,
		ctx:    ctx,
		agent:  agent,
		sd:     sd,
	}

	go self.Update()

	return self
}
