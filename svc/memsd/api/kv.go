package memsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/svc/memsd/proto"
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

	self.remoteCall(&sdproto.SetValueREQ{
		Key:   key,
		Value: raw,
	}, func(ack *sdproto.SetValueACK, err error) {

		if err != nil {
			retErr = err
		} else {
			retErr = codeToError(ack.Code)
		}
	})

	return
}

func (self *memDiscovery) GetValue(key string, valuePtr interface{}) error {

	data, ok := self.getKVCache(key)

	if !ok {
		return ErrValueNotExists
	}

	return discovery.BytesToAny(data, valuePtr)
}

func (self *memDiscovery) GetRawValue(key string) (retData []byte, retErr error) {

	self.remoteCall(&sdproto.GetValueREQ{
		Key: key,
	}, func(ack *sdproto.GetValueACK, err error) {
		if err != nil {
			retErr = err
		} else {
			retData = ack.Value
			retErr = codeToError(ack.Code)
		}

	})

	return
}

func (self *memDiscovery) DeleteValue(key string) (ret error) {

	self.remoteCall(&sdproto.DeleteValueREQ{
		Key: key,
	}, func(ack *sdproto.DeleteValueACK, err error) {
		if err != nil {
			ret = err
		} else {
			ret = codeToError(ack.Code)
		}

	})

	return
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
	self.remoteCall(&sdproto.ClearKeyREQ{}, func(ack *sdproto.ClearKeyACK, err error) {})
}
