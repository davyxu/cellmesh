package consulsd

import (
	"errors"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
	"time"
)

func (self *consulDiscovery) SetValue(key string, dataPtr interface{}) error {

	raw, err := AnyToBytes(dataPtr)
	if err != nil {
		return err
	}

	_, err = self.client.KV().Put(&api.KVPair{
		Key:   key,
		Value: raw,
	}, nil)

	return err
}

func getOption(opts []interface{}) (ret discovery.Option) {
	for _, raw := range opts {

		if opt, ok := raw.(discovery.Option); ok {
			ret = opt
		}
	}

	return
}

func (self *consulDiscovery) GetValue(key string, valuePtr interface{}, opts ...interface{}) error {

	if getOption(opts).NoCache {

		value, err := self.directGetValue(key)

		if err != nil {
			return err
		}

		return discovery.BytesToAny(value, valuePtr)

	} else {

		data, err := self.GetRawValue(key)
		if err != nil {
			return err
		}

		return discovery.BytesToAny(data, valuePtr)
	}

}

func (self *consulDiscovery) GetRawValue(key string) ([]byte, error) {
	cv, err := self.getKV(key, func(value *cacheValue) (thisErr error) {
		value.value, thisErr = self.directGetValue(key)
		return
	})

	if err != nil {
		return nil, err
	}

	return cv.value, nil
}

func (self *consulDiscovery) GetRawValueList(key string) (ret []discovery.ValueMeta, err error) {
	cv, err := self.getKV(key, func(value *cacheValue) (thisErr error) {
		value.valueList, thisErr = self.directGetValueList(key)
		return
	})

	if err != nil {
		return nil, err
	}

	for _, v := range cv.valueList {
		ret = append(ret, discovery.ValueMeta{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	return
}

var (
	ErrValueNotExists = errors.New("value not exists")
)

func (self *consulDiscovery) directGetValue(key string) ([]byte, error) {
	kvPair, _, err := self.client.KV().Get(key, nil)

	if err != nil {
		return nil, err
	}

	// 值不存在
	if kvPair == nil {
		return nil, ErrValueNotExists
	}

	return kvPair.Value, nil
}

func (self *consulDiscovery) directGetValueList(key string) (api.KVPairs, error) {
	kvPairs, _, err := self.client.KV().List(key, nil)

	if err != nil {
		return nil, err
	}

	return kvPairs, nil
}

type cacheValue struct {
	value     []byte
	valueList api.KVPairs

	time time.Time
}

func (self *consulDiscovery) getKV(key string, callback func(*cacheValue) error) (*cacheValue, error) {
	now := time.Now()

	raw, ok := self.metaByKey.Load(key)

	var cv *cacheValue
	if ok {
		cv = raw.(*cacheValue)

		if now.Sub(cv.time) < self.config.KVCacheDuration {
			return cv, nil
		}
	} else {
		cv = &cacheValue{} // 不存在，创建新的
	}

	err := callback(cv)
	if err != nil {
		return nil, err
	}

	cv.time = now

	// 刷新缓冲
	self.metaByKey.Store(key, cv)

	return cv, err
}
