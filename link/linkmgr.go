package link

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/ulog"
)

func getPeerDescString(peer cellnet.Peer) string {

	desc := GetPeerDesc(peer)
	if desc != nil {
		return fmt.Sprintf("%s", desc.Address())
	}

	return ""
}

// 内部逻辑使用
func markLink(ses cellnet.Session, svcid, svcName string) {
	ctxSet := ses.(cellnet.ContextSet)
	ctxSet.SetContext(SesContextKey_LinkSvcID, svcid)
	ctxSet.SetContext(SesContextKey_LinkSvcName, svcName)

	ulog.WithColorName("green").Infof("Add service link: %s %s", svcid, getPeerDescString(ses.Peer()))
}

// 取得远程会话的ID
func GetLinkSvcID(ses cellnet.Session) string {
	if ses == nil {
		return ""
	}

	if raw, ok := ses.(cellnet.ContextSet).GetContext(SesContextKey_LinkSvcID); ok {
		return raw.(string)
	}

	return ""
}

// 取得远程会话的服务名
func GetLinkSvcName(ses cellnet.Session) string {
	if ses == nil {
		return ""
	}

	if raw, ok := ses.(cellnet.ContextSet).GetContext(SesContextKey_LinkSvcName); ok {
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

	VisitPeer(func(peer cellnet.Peer) bool {

		linkSes := GetPeerLink(peer)
		if linkSes != nil {
			return callback(linkSes)
		}

		return true
	})

}

func printAllLink() {
	VisitLink(func(ses cellnet.Session) bool {
		ulog.Debugf("link: ", GetLinkSvcID(ses))
		return true
	})
}

func OneLink(svcName string) (ret cellnet.Session) {
	VisitLink(func(ses cellnet.Session) bool {
		if GetLinkSvcName(ses) == svcName {
			ret = ses
			return false
		}

		return true
	})

	return
}
