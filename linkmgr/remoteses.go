package linkmgr

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellnet"
	"sync"
)

type NotifyFunc func(ses cellnet.Session)

var (
	sesBySvcID      = map[string]cellnet.Session{}
	sesBySvcIDGuard sync.RWMutex
	removeNotify    NotifyFunc
)

func AddRemoteSession(ses cellnet.Session, svcid string) {

	sesBySvcIDGuard.Lock()
	ses.(cellnet.ContextSet).SetContext(cellmesh.RemoteSesKey_RemoteSvcID, svcid)
	sesBySvcID[svcid] = ses
	sesBySvcIDGuard.Unlock()

	log.SetColor("green").Infof("remote service added: '%s' sid: %d", svcid, ses.ID())
}

func RemoveRemoteSession(ses cellnet.Session) {

	if ses == nil {
		return
	}

	svcID := GetRemoteSessionSvcID(ses)
	if svcID != "" {

		if removeNotify != nil {
			removeNotify(ses)
		}

		sesBySvcIDGuard.Lock()
		delete(sesBySvcID, svcID)
		sesBySvcIDGuard.Unlock()

		log.SetColor("yellow").Infof("remote service removed '%s' sid: %d", svcID, ses.ID())
	} else {
		log.SetColor("yellow").Infof("remote service removed sid: %d, context lost", ses.ID())
	}
}

// 设置服务的通知
func SetRemoteSessionNotify(mode string, callback NotifyFunc) {

	switch mode {
	case "remove":
		removeNotify = callback
	default:
		panic("unknown notify mode")
	}
}

// 取得其他服务器的会话对应的上下文
func GetRemoteSessionSvcID(ses cellnet.Session) string {
	if ses == nil {
		return ""
	}

	if raw, ok := ses.(cellnet.ContextSet).GetContext(cellmesh.RemoteSesKey_RemoteSvcID); ok {
		return raw.(string)
	}

	return ""
}

// 根据svcid获取远程服务的会话
func GetRemoteSession(svcid string) cellnet.Session {
	sesBySvcIDGuard.RLock()
	defer sesBySvcIDGuard.RUnlock()

	if ses, ok := sesBySvcID[svcid]; ok {

		return ses
	}

	return nil
}

// 遍历远程服务(已经连接到本进程)
func VisitRemoteSession(callback func(ses cellnet.Session) bool) {
	sesBySvcIDGuard.RLock()

	for _, ses := range sesBySvcID {

		if !callback(ses) {
			break
		}
	}

	sesBySvcIDGuard.RUnlock()
}
