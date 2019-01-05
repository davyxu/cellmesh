package memsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/memsd/proto"
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

func (self *memDiscovery) getKVCache(key string) (value []byte, ok bool) {
	self.kvCacheGuard.RLock()
	defer self.kvCacheGuard.RUnlock()
	value, ok = self.kvCache[key]
	return
}

func (self *memDiscovery) updateKVCache(key string, value []byte) {
	self.kvCacheGuard.Lock()
	self.kvCache[key] = value
	self.kvCacheGuard.Unlock()
}

func (self *memDiscovery) deleteKVCache(key string) {
	self.kvCacheGuard.Lock()
	delete(self.kvCache, key)
	self.kvCacheGuard.Unlock()
}

func (self *memDiscovery) SetValue(key string, dataPtr interface{}, optList ...interface{}) (retErr error) {

	raw, err := discovery.AnyToBytes(dataPtr, getOpt(optList...).PrettyPrint)
	if err != nil {
		return err
	}

	if len(raw) > MaxValueSize {
		return ErrValueTooLarge
	}

	callErr := self.remoteCall(&proto.SetValueREQ{
		Key:   key,
		Value: raw,
	}, func(ack *proto.SetValueACK) {
		retErr = codeToError(ack.Code)
	})

	if retErr != nil {
		return
	}

	retErr = callErr

	return nil
}

func (self *memDiscovery) GetValue(key string, valuePtr interface{}) error {

	data, ok := self.getKVCache(key)

	if !ok {
		return ErrValueNotExists
	}

	return discovery.BytesToAny(data, valuePtr)
}

func (self *memDiscovery) GetRawValue(key string) (retData []byte, retErr error) {

	callErr := self.remoteCall(&proto.GetValueREQ{
		Key: key,
	}, func(ack *proto.GetValueACK) {
		retData = ack.Value
		retErr = codeToError(ack.Code)
	})

	if retErr != nil {
		return
	}

	retErr = callErr

	return
}

func (self *memDiscovery) DeleteValue(key string) (ret error) {

	callErr := self.remoteCall(&proto.DeleteValueREQ{
		Key: key,
	}, func(ack *proto.DeleteValueACK) {
		ret = codeToError(ack.Code)
	})

	if ret != nil {
		return ret
	}

	return callErr
}

func (self *memDiscovery) GetRawValueList(prefix string) (ret []discovery.ValueMeta) {

	self.kvCacheGuard.RLock()

	for key, value := range self.kvCache {

		if strings.HasPrefix(key, prefix) {
			ret = append(ret, discovery.ValueMeta{
				Key:   key,
				Value: value,
			})
		}

	}

	self.kvCacheGuard.RUnlock()

	return
}

func (self *memDiscovery) ClearKey() {
	self.remoteCall(&proto.ClearKeyREQ{}, func(ack *proto.ClearKeyACK) {})
}
