package service

import (
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

	log.SetColor("green").Debugf("remote service added: '%s'", desc.ID)
}

func RemoveRemoteService(ses cellnet.Session) {
	if raw, ok := ses.(cellnet.ContextSet).GetContext("desc"); ok {

		desc := raw.(*discovery.ServiceDesc)

		connBySvcNameGuard.Lock()
		delete(connBySvcID, desc.ID)
		connBySvcNameGuard.Unlock()

		log.SetColor("yellow").Debugf("remote service removed '%s'", desc.ID)
	}
}

// 获取Session绑定的远程连接描述
func SessionToDesc(ses cellnet.Session) *discovery.ServiceDesc {

	if raw, ok := ses.(cellnet.ContextSet).GetContext("desc"); ok {
		return raw.(*discovery.ServiceDesc)
	}

	return nil
}

func GetRemoteService(svcid string) cellnet.Session {
	connBySvcNameGuard.RLock()
	defer connBySvcNameGuard.RUnlock()

	if ses, ok := connBySvcID[svcid]; ok {

		return ses
	}

	return nil
}

// 遍历远程服务
func VisitRemoteService(callback func(ses cellnet.Session, desc *discovery.ServiceDesc) bool) {
	connBySvcNameGuard.RLock()

	for _, ses := range connBySvcID {

		sd := SessionToDesc(ses)

		if !callback(ses, sd) {
			break
		}
	}

	connBySvcNameGuard.RUnlock()
}
