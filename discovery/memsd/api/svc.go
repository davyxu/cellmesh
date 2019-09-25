package memsd

import (
	"encoding/json"
	"errors"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/memsd/model"
	"github.com/davyxu/cellmesh/discovery/memsd/proto"
)

func (self *memDiscovery) Register(svc *discovery.ServiceDesc) (retErr error) {

	if svc.Name == "" {
		return errors.New("expect svc name")
	}

	if svc.ID == "" {
		return errors.New("expect svc id")
	}

	data, err := json.Marshal(svc)
	if err != nil {
		return err
	}

	callErr := self.remoteCall(&sdproto.SetValueREQ{
		Key:     model.ServiceKeyPrefix + svc.ID,
		Value:   data,
		SvcName: svc.Name,
	}, func(ack *sdproto.SetValueACK) {
		retErr = codeToError(ack.Code)
	})

	if retErr != nil {
		return
	}

	return callErr
}

func (self *memDiscovery) Deregister(svcid string) error {

	return self.DeleteValue(model.ServiceKeyPrefix + svcid)
}

func (self *memDiscovery) Query(name string) (ret []*discovery.ServiceDesc) {

	self.svcCacheGuard.RLock()
	defer self.svcCacheGuard.RUnlock()

	return self.svcCache[name]
}

func (self *memDiscovery) QueryAll() (ret []*discovery.ServiceDesc) {

	self.svcCacheGuard.RLock()
	defer self.svcCacheGuard.RUnlock()

	for _, list := range self.svcCache {
		ret = append(ret, list...)
	}

	return
}

func (self *memDiscovery) ClearService() {
	self.remoteCall(&sdproto.ClearSvcREQ{}, func(ack *sdproto.ClearSvcACK) {})
}

func (self *memDiscovery) updateSvcCache(svcName string, value []byte) {
	self.svcCacheGuard.Lock()

	list := self.svcCache[svcName]

	var desc discovery.ServiceDesc
	err := json.Unmarshal(value, &desc)
	if err != nil {
		log.Errorf("ServiceDesc unmarshal failed, %s", err)
		self.svcCacheGuard.Unlock()
		return
	}

	var notfound = true
	for index, svc := range list {
		if svc.ID == desc.ID {
			list[index] = &desc
			notfound = false
			break
		}
	}

	if notfound {
		list = append(list, &desc)
	}

	self.svcCache[svcName] = list
	self.svcCacheGuard.Unlock()

	self.triggerNotify("add", &desc)
}

func (self *memDiscovery) deleteSvcCache(svcid, svcName string) {

	list := self.svcCache[svcName]

	for index, svc := range list {
		if svc.ID == svcid {
			list = append(list[:index], list[index+1:]...)
			break
		}
	}

	self.svcCache[svcName] = list
}
