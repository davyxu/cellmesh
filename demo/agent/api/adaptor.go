package api

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
)

// 传入用户处理网关消息回调,返回消息源回调
func HandleBackendMessage(userHandler func(ev cellnet.Event, cid proto.ClientID)) func(ev cellnet.Event) {

	return func(ev cellnet.Event) {

		var cid proto.ClientID
		if err := service.GetPassThrough(ev, &cid.ID, &cid.SvcID); err != nil {
			log.Errorf("service.GetPassThrough %s", err)
		} else {
			userHandler(ev, cid)
		}

	}
}
