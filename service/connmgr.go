package service

import (
	"errors"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"sync"
)

var (
	connBySvcID        = map[string]cellnet.Session{}
	connBySvcNameGuard sync.RWMutex
)

func AddRemoteService(ses cellnet.Session, desc *discovery.ServiceDesc) {

	connBySvcNameGuard.Lock()
	ses.(cellnet.ContextSet).SetContext("desc", desc)
	connBySvcID[desc.ID] = ses
	connBySvcNameGuard.Unlock()

	log.SetColor("green").Debugf("add connection: '%s'", desc.ID)
}

func GetRemoteService(svcid string) cellnet.Session {
	connBySvcNameGuard.RLock()
	defer connBySvcNameGuard.RUnlock()

	if ses, ok := connBySvcID[svcid]; ok {

		return ses
	}

	return nil
}

func ServiceSessionToSD(ses cellnet.Session) *discovery.ServiceDesc {

	if raw, ok := ses.(cellnet.ContextSet).GetContext("desc"); ok {
		return raw.(*discovery.ServiceDesc)
	}

	return nil
}

func RemoveRemoteService(ses cellnet.Session) {
	if raw, ok := ses.(cellnet.ContextSet).GetContext("desc"); ok {

		desc := raw.(*discovery.ServiceDesc)

		connBySvcNameGuard.Lock()
		delete(connBySvcID, desc.ID)
		connBySvcNameGuard.Unlock()

		log.SetColor("yellow").Debugf("connection remove '%s'", desc.ID)
	}
}

func VisitRemoteService(callback func(ses cellnet.Session, desc *discovery.ServiceDesc) bool) {
	connBySvcNameGuard.RLock()

	for _, ses := range connBySvcID {

		sd := ServiceSessionToSD(ses)

		if !callback(ses, sd) {
			break
		}
	}

	connBySvcNameGuard.RUnlock()
}

func selectStrategy(descList []*discovery.ServiceDesc) *discovery.ServiceDesc {

	if len(descList) == 0 {
		return nil
	}

	return descList[0]
}

func QueryServiceAddress(serviceName string) (*discovery.ServiceDesc, error) {
	descList, err := discovery.Default.Query(serviceName)
	if err != nil {
		return nil, err
	}

	desc := selectStrategy(descList)

	if desc == nil {
		return nil, errors.New("target not reachable:" + serviceName)
	}

	return desc, nil
}
