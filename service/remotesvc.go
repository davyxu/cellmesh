package service

import (
	"github.com/davyxu/cellnet"
	"sync"
)

type RemoteServiceContext struct {
	Name  string
	SvcID string
}

type NotifyFunc func(ctx *RemoteServiceContext, ses cellnet.Session)

var (
	connBySvcID        = map[string]cellnet.Session{}
	connBySvcNameGuard sync.RWMutex
	removeNotify       NotifyFunc
)

func AddRemoteService(ses cellnet.Session, svcid, name string) {

	connBySvcNameGuard.Lock()
	ses.(cellnet.ContextSet).SetContext("ctx", &RemoteServiceContext{Name: name, SvcID: svcid})
	connBySvcID[svcid] = ses
	connBySvcNameGuard.Unlock()

	log.SetColor("green").Debugf("remote service added: '%s'", svcid)
}

func RemoveRemoteService(ses cellnet.Session) {

	ctx := SessionToContext(ses)
	if ctx != nil {

		if removeNotify != nil {
			removeNotify(ctx, ses)
		}

		connBySvcNameGuard.Lock()
		delete(connBySvcID, ctx.SvcID)
		connBySvcNameGuard.Unlock()

		log.SetColor("yellow").Debugf("remote service removed '%s'", ctx.SvcID)
	}
}

// 设置服务的通知
func SetRemoteServiceNotify(mode string, callback NotifyFunc) {

	switch mode {
	case "remove":
		removeNotify = callback
	default:
		panic("unknown notify mode")
	}
}

// 取得其他服务器的会话对应的上下文
func SessionToContext(ses cellnet.Session) *RemoteServiceContext {
	if raw, ok := ses.(cellnet.ContextSet).GetContext("ctx"); ok {
		return raw.(*RemoteServiceContext)
	}

	return nil
}

// 根据svcid获取远程服务的会话
func GetRemoteService(svcid string) cellnet.Session {
	connBySvcNameGuard.RLock()
	defer connBySvcNameGuard.RUnlock()

	if ses, ok := connBySvcID[svcid]; ok {

		return ses
	}

	return nil
}

// 遍历远程服务(已经连接到本进程)
func VisitRemoteService(callback func(ses cellnet.Session, ctx *RemoteServiceContext) bool) {
	connBySvcNameGuard.RLock()

	for _, ses := range connBySvcID {

		if !callback(ses, SessionToContext(ses)) {
			break
		}
	}

	connBySvcNameGuard.RUnlock()
}
