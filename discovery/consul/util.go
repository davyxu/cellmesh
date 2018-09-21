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
			check.Output == MakeHealthWords(entry.Service.ID) {
			return true
		}
	}

	// 非内建建立的服务， 比如redis
	if len(entry.Checks) > 0 && entry.Checks[0].Status == "passing" {
		return true
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

var (
	PrettyMarshalJson bool
)

func AnyToBytes(data interface{}) ([]byte, error) {

	switch v := data.(type) {
	case int, int32, int64, uint32, uint64, float32, float64, bool:
		return []byte(fmt.Sprint(data)), nil
	case string:
		return []byte(v), nil

	default:
		if PrettyMarshalJson {
			raw, err := json.MarshalIndent(data, "", "\t")
			if err != nil {
				return nil, err
			}

			return raw, nil
		} else {
			raw, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			return raw, nil
		}
	}
}
