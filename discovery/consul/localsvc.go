package consulsd

import (
	"context"
	"github.com/hashicorp/consul/api"
	"time"
)

// 本地服务更新TTL
type localService struct {
	ID     string
	Cancel context.CancelFunc

	ctx context.Context

	agent *api.Agent
}

func (self *localService) Update() {

	//log.Debugf("UpdateTTL id: %s begin", self.ID)

	for {

		select {
		case <-self.ctx.Done():
			return
		default:

			//log.Debugf("UpdateTTL id: %s", self.ID)

			self.agent.UpdateTTL(self.ID, "svc ready", "pass")

			time.Sleep(ServiceTTL)
		}
	}

	//log.Debugf("UpdateTTL id: %s end", self.ID)
}

func (self *localService) Stop() {
	self.Cancel()
}

func newLocalService(id string, agent *api.Agent) *localService {

	ctx, cancel := context.WithCancel(context.Background())

	self := &localService{
		ID:     id,
		Cancel: cancel,
		ctx:    ctx,
		agent:  agent,
	}

	go self.Update()

	return self
}
