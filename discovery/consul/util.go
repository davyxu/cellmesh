package consulsd

import (
	"encoding/json"
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
)

// check包括：node的check和service check
func isMeshServiceHealth(entry *api.ServiceEntry) bool {

	for _, check := range entry.Checks {
		if check.ServiceID == entry.Service.ID &&
			check.Output == makeHealthWords(entry.Service.ID) {
			return true
		}
	}

	return false
}

func consulSvcToService(s *api.ServiceEntry) *discovery.ServiceDesc {

	return &discovery.ServiceDesc{
		Name: s.Service.Service,
		ID:   s.Service.ID,
		Host: s.Service.Address,
		Port: s.Service.Port,
		Tags: s.Service.Tags,
		Meta: s.Service.Meta,
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

func AnyToBytes(data interface{}) ([]byte, error) {

	switch v := data.(type) {
	case int, int32, int64, uint32, uint64, float32, float64, bool:
		return []byte(fmt.Sprint(data)), nil
	case string:
		return []byte(v), nil

	default:
		raw, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		return raw, nil
	}
}
