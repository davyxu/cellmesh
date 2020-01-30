package agentapi

import (
	"github.com/davyxu/cellmesh/link"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellmesh/svc/agent/backend"
	"github.com/davyxu/cellnet"
)

// 传入用户处理网关消息回调,返回消息源回调
func HandleBackendMessage(userHandler func(ev cellnet.Event, cid proto.ClientID)) func(ev cellnet.Event) {

	return func(incomingEv cellnet.Event) {

		switch ev := incomingEv.(type) {
		case *backend.RecvMsgEvent:

			var cid proto.ClientID
			cid.ID = ev.ClientID

			if agentCtx := link.GetPeerDesc(ev.Session().Peer()); agentCtx != nil {
				cid.SvcID = agentCtx.ID
			}

			userHandler(incomingEv, cid)
		}
	}
}
