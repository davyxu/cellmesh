package consulsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
)

// check包括：node的check和service check
func isServiceHealth(entry *api.ServiceEntry) bool {

	for _, check := range entry.Checks {
		if check.Status != "passing" {
			return false
		}
	}

	return true
}

func consulSvcToService(s *api.AgentService) *discovery.ServiceDesc {

	return &discovery.ServiceDesc{
		Name: s.Service,
		ID:   s.ID,
		Host: s.Address,
		Port: s.Port,
		Tags: s.Tags,
		Meta: s.Meta,
	}
}

func existsInServiceList(svclist []*discovery.ServiceDesc, id string) bool {
	for _, svc := range svclist {

		if svc.ID == id {
			return true
		}
	}

	return false
}
