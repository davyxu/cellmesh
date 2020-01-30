package link

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellnet"
)

// 内部逻辑使用
func markLink(ses cellnet.Session, svcid, svcName string) {
	ctxSet := ses.(cellnet.ContextSet)
	ctxSet.SetContext(cellmesh.SesContextKey_LinkSvcID, svcid)
	ctxSet.SetContext(cellmesh.SesContextKey_LinkSvcName, svcName)

	log.SetColor("green").Infof("add remote service : '%s' sid: %d", svcid, ses.ID())
}

// 取得远程会话的ID
func GetLinkSvcID(ses cellnet.Session) string {
	if ses == nil {
		return ""
	}

	if raw, ok := ses.(cellnet.ContextSet).GetContext(cellmesh.SesContextKey_LinkSvcID); ok {
		return raw.(string)
	}

	return ""
}

// 取得远程会话的服务名
func GetLinkSvcName(ses cellnet.Session) string {
	if ses == nil {
		return ""
	}

	if raw, ok := ses.(cellnet.ContextSet).GetContext(cellmesh.SesContextKey_LinkSvcName); ok {
		return raw.(string)
	}

	return ""
}

// 根据svcid获取远程服务的会话
func GetLink(svcid string) (ret cellnet.Session) {

	VisitLink(func(ses cellnet.Session) bool {

		if GetLinkSvcID(ses) == svcid {
			ret = ses
			return false
		}

		return true
	})
	return
}

// 遍历远程服务(已经连接到本进程)
func VisitLink(callback func(ses cellnet.Session) bool) {

	visitContinue := true
	VisitPeer(func(peer cellnet.Peer) bool {

		// 注意, 已经断开的session的Connector是没有session的

		if sesmgr, ok := peer.(cellnet.SessionAccessor); ok {
			sesmgr.VisitSession(func(ses cellnet.Session) bool {

				visitContinue = callback(ses)

				return visitContinue
			})

		} else {
			log.Errorf("peer not support 'cellnet.SessionAccessor', %s", peer.TypeName())

		}

		return visitContinue
	})

}
