package consulsd

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"reflect"
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

func (self *consulDiscovery) GetValue(key string, dataPtr interface{}) error {

	vdata := reflect.Indirect(reflect.ValueOf(dataPtr))
	if vdata.Kind() == reflect.Slice {

		cv, err := self.getKV(key, func(value *cacheValue) (thisErr error) {
			value.valueList, thisErr = self.directGetValueList(key)
			return
		})

		if err != nil {
			return err
		}

		if err := pairsToSlice(cv.valueList, dataPtr); err != nil {
			return err
		}

		return nil

	} else {
		cv, err := self.getKV(key, func(value *cacheValue) (thisErr error) {
			value.value, thisErr = self.directGetValue(key)
			return
		})

		if err != nil {
			return err
		}

		// TODO 本地缓存，及轮询/watch更新
		return BytesToAny(cv.value, dataPtr)
	}
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

		if now.Sub(cv.time) < self.cacheDuration {
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
