package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"sync"
)

type remoteContext struct {
	name string
	id   string
}

var (
	connBySvcID        = map[string]cellnet.Session{}
	connBySvcNameGuard sync.RWMutex
)

func AddRemoteService(ses cellnet.Session, svcid, name string) {

	connBySvcNameGuard.Lock()
	ses.(cellnet.ContextSet).SetContext("ctx", &remoteContext{name: name, id: svcid})
	connBySvcID[svcid] = ses
	connBySvcNameGuard.Unlock()

	log.SetColor("green").Debugf("remote service added: '%s'", svcid)
}

func RemoveRemoteService(ses cellnet.Session) {
	if raw, ok := ses.(cellnet.ContextSet).GetContext("ctx"); ok {

		ctx := raw.(*remoteContext)

		connBySvcNameGuard.Lock()
		delete(connBySvcID, ctx.id)
		connBySvcNameGuard.Unlock()

		log.SetColor("yellow").Debugf("remote service removed '%s'", ctx.id)
	}
}

// 获取Session绑定的远程连接描述
func SessionToDesc(ses cellnet.Session) *discovery.ServiceDesc {
	if ses == nil {
		return nil
	}

	if raw, ok := ses.(cellnet.ContextSet).GetContext("ctx"); ok {
		ctx := raw.(*remoteContext)

		// 要取新鲜的
		descList := DiscoveryService(LinkRules, ctx.name)
		for _, desc := range descList {
			if desc.ID == ctx.id {
				return desc
			}
		}
	}

	return nil
}

// 根据id获取远程服务的会话
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

// 获得一个远程服务会话的外网地址
func GetRemoteServiceWANAddress(ses cellnet.Session) string {
	desc := SessionToDesc(ses)

	if desc == nil {
		return ""
	}

	return desc.GetMeta("WANAddress")
}
