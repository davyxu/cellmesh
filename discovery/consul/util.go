package consulsd

import (
	"encoding/json"
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
	"reflect"
	"strconv"
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

func BytesToAny(data []byte, dataPtr interface{}) error {

	switch ret := dataPtr.(type) {
	case *int:
		v, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		*ret = int(v)
		return nil
	case *float32:
		v, err := strconv.ParseFloat(string(data), 32)
		if err != nil {
			return err
		}
		*ret = float32(v)
		return nil
	case *float64:
		v, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		*ret = float64(v)
		return nil
	case *bool:
		v, err := strconv.ParseBool(string(data))
		if err != nil {
			return err
		}
		*ret = v
		return nil
	case *string:
		*ret = string(data)
		return nil
	default:
		return json.Unmarshal(data, dataPtr)
	}
}

func pairsToSlice(pairs api.KVPairs, dataPtr interface{}) error {

	vdata := reflect.Indirect(reflect.ValueOf(dataPtr))

	elementCount := len(pairs)

	slice := reflect.MakeSlice(vdata.Type(), elementCount, elementCount)

	for i := 0; i < elementCount; i++ {

		sliceValue := reflect.New(slice.Type().Elem())

		err := BytesToAny(pairs[i].Value, sliceValue.Interface())

		if err != nil {
			return err
		}

		slice.Index(i).Set(sliceValue.Elem())
	}

	vdata.Set(slice)

	return nil

}
