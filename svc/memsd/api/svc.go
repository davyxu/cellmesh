package memsd

import (
	"encoding/json"
	"errors"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/svc/memsd/model"
	"github.com/davyxu/cellmesh/svc/memsd/proto"
	"github.com/davyxu/cellnet/util"
	"time"
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

	self.remoteCall(&sdproto.SetValueREQ{
		Key:     model.ServiceKeyPrefix + svc.ID,
		Value:   data,
		SvcName: svc.Name,
	}, func(ack *sdproto.SetValueACK, err error) {
		if err != nil {
			retErr = err
		} else {
			retErr = codeToError(ack.Code)
		}

	})

	// 确保信息已经同步到本地
	for {
		descList := self.Query(svc.Name)

		if discovery.DescExistsByID(svc.ID, descList) {
			break
		}

		time.Sleep(time.Millisecond * 100)
	}

	return
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
	self.remoteCall(&sdproto.ClearSvcREQ{}, func(ack *sdproto.ClearSvcACK, err error) {
		if err != nil {
			log.Errorln(err)
		}
	})
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

		self.triggerNotify("add", &desc)

		list = append(list, &desc)
	} else {
		self.triggerNotify("mod", &desc)
	}

	self.svcCache[svcName] = list
	self.svcCacheGuard.Unlock()

}

type notifyContext struct {
	stack string
}

func (self *memDiscovery) RegisterNotify() (ret chan *discovery.NotifyContext) {
	ret = make(chan *discovery.NotifyContext, 100)

	self.notifyMap.Store(ret, &notifyContext{
		stack: util.StackToString(5),
	})

	return
}

func (self *memDiscovery) DeregisterNotify(c chan struct{}) {

	self.notifyMap.Store(c, nil)
}

func (self *memDiscovery) triggerNotify(mode string, desc *discovery.ServiceDesc) {

	self.notifyMap.Range(func(key, value interface{}) bool {

		if value == nil {
			return true
		}

		ctx := value.(*notifyContext)

		c := key.(chan *discovery.NotifyContext)

		select {
		case c <- &discovery.NotifyContext{
			Mode: mode,
			Desc: desc,
		}:
		case <-time.After(time.Second * 10):
			// 接收通知阻塞太久，或者没有释放侦听的channel
			log.Errorf("notify(%s) timeout, not free? regstack: %s, desc: %s", mode, ctx.stack, desc.String())
		}

		return true
	})

}

func (self *memDiscovery) deleteSvcCache(svcid, svcName string) {

	list := self.svcCache[svcName]

	var svcToRemove []*discovery.ServiceDesc
	for index, svc := range list {
		if svc.ID == svcid {

			svcToRemove = append(svcToRemove, svc)
			list = append(list[:index], list[index+1:]...)
			break
		}
	}

	self.svcCache[svcName] = list

	for _, svc := range svcToRemove {
		self.triggerNotify("del", svc)
	}
}
