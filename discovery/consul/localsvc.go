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

	checkerFunc discovery.CheckerFunc

	lastOutput string
}

func (self *localService) Update() {

	for {

		select {
		case <-self.ctx.Done():
			return
		default:

			var output string
			var status string

			if self.checkerFunc != nil {

				output, status = self.checkerFunc()
			} else {
				output = self.Desc.ID
				status = "pass"
			}

			// 注意，consul这里有bug https://github.com/hashicorp/consul/issues/1057
			// 只有在status变化和有服务加入时，status才能及时更新，但是output依然不能及时更新
			self.agent.UpdateTTL(self.Desc.ID, output, status)

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
