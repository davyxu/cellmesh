package link

import (
	"github.com/davyxu/cellmesh"
	meshutil "github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"sync"
)

var (
	// 连接建立
	OnLinkAdd meshutil.EventFuncSet

	// 连接关闭
	OnLinkRemove meshutil.EventFuncSet
)

var (
	linkBySvcID      = map[string]cellnet.Session{}
	linkBySvcIDGuard sync.RWMutex
)

func AddRemoteLink(ses cellnet.Session, svcid, svcName string) {

	linkBySvcIDGuard.Lock()
	ctxSet := ses.(cellnet.ContextSet)
	ctxSet.SetContext(cellmesh.SesContextKey_LinkSvcID, svcid)
	ctxSet.SetContext(cellmesh.SesContextKey_LinkSvcName, svcName)
	linkBySvcID[svcid] = ses
	linkBySvcIDGuard.Unlock()

	log.SetColor("green").Infof("add remote service : '%s' sid: %d", svcid, ses.ID())

	OnLinkAdd.Invoke(ses)
}

func RemoveRemoteLink(ses cellnet.Session) {

	if ses == nil {
		return
	}

	svcID := GetRemoteLinkSvcID(ses)
	if svcID != "" {

		OnLinkRemove.Invoke(ses)

		linkBySvcIDGuard.Lock()
		delete(linkBySvcID, svcID)
		linkBySvcIDGuard.Unlock()

		log.SetColor("yellow").Infof("remove remote service '%s' sid: %d", svcID, ses.ID())
	} else {
		log.SetColor("yellow").Infof("remove service sid: %d, context lost", ses.ID())
	}
}

// 取得远程会话的ID
func GetRemoteLinkSvcID(ses cellnet.Session) string {
	if ses == nil {
		return ""
	}

	if raw, ok := ses.(cellnet.ContextSet).GetContext(cellmesh.SesContextKey_LinkSvcID); ok {
		return raw.(string)
	}

	return ""
}

// 取得远程会话的服务名
func GetRemoteLinkSvcName(ses cellnet.Session) string {
	if ses == nil {
		return ""
	}

	if raw, ok := ses.(cellnet.ContextSet).GetContext(cellmesh.SesContextKey_LinkSvcName); ok {
		return raw.(string)
	}

	return ""
}

// 根据svcid获取远程服务的会话
func GetRemoteLink(svcid string) cellnet.Session {
	linkBySvcIDGuard.RLock()
	defer linkBySvcIDGuard.RUnlock()

	if ses, ok := linkBySvcID[svcid]; ok {

		return ses
	}

	return nil
}

// 遍历远程服务(已经连接到本进程)
func VisitRemoteSession(callback func(ses cellnet.Session) bool) {
	linkBySvcIDGuard.RLock()

	for _, ses := range linkBySvcID {

		if !callback(ses) {
			break
		}
	}

	linkBySvcIDGuard.RUnlock()
}
