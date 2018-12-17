package consulsd

import (
	"errors"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
	"strings"
)

type Option struct {
	PrettyPrint bool
}

func getOpt(optList ...interface{}) Option {

	for _, opt := range optList {

		switch raw := opt.(type) {
		case Option:
			return raw
		}
	}

	return Option{}
}

func (self *consulDiscovery) SetValue(key string, dataPtr interface{}, optList ...interface{}) error {

	raw, err := AnyToBytes(dataPtr, getOpt(optList...).PrettyPrint)
	if err != nil {
		return err
	}

	_, err = self.client.KV().Put(&api.KVPair{
		Key:   key,
		Value: raw,
	}, nil)

	return err
}

func (self *consulDiscovery) GetValue(key string, valuePtr interface{}) error {

	data, err := self.GetRawValue(key)
	if err != nil {
		return err
	}

	return discovery.BytesToAny(data, valuePtr)
}

func (self *consulDiscovery) GetRawValue(key string) ([]byte, error) {

	if raw, ok := self.kvCache.Load(key); ok {

		meta := raw.(*KVMeta)

		return meta.Value(), nil
	} else {

		// cache中没找到直接获取
		kvpair, _, err := self.client.KV().Get(key, nil)
		if err != nil {
			return nil, err
		}

		if kvpair == nil {
			return nil, ErrValueNotExists
		}

		return kvpair.Value, nil
	}
}

func (self *consulDiscovery) GetRawValueList(prefix string) (ret []discovery.ValueMeta, err error) {

	self.kvCache.Range(func(rawKey, rawValue interface{}) bool {

		key := rawKey.(string)
		value := rawValue.(*KVMeta)

		if strings.HasPrefix(key, prefix) {
			ret = append(ret, discovery.ValueMeta{
				Key:   key,
				Value: value.Value(),
			})
		}

		return true

	})

	if len(ret) == 0 {
		kvpairs, _, err := self.client.KV().List(prefix, nil)
		if err != nil {
			return nil, err
		}

		for _, kv := range kvpairs {
			ret = append(ret, discovery.ValueMeta{
				Key:   kv.Key,
				Value: kv.Value,
			})
		}
	}

	return
}

func (self *consulDiscovery) DeleteValue(key string) error {

	_, err := self.client.KV().Delete(key, nil)
	if err != nil {
		return err
	}

	return nil
}

var (
	ErrValueNotExists = errors.New("value not exists")
)
