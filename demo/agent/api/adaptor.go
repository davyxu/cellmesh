package api

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
)

// 传入用户处理网关消息回调,返回消息源回调
func HandleBackendMessage(userHandler func(ev cellnet.Event, cid proto.ClientID)) func(ev cellnet.Event) {

	return func(ev cellnet.Event) {

		if cid, ok := service.GetPassThrough(ev).(*proto.ClientID); ok {
			userHandler(ev, *cid)
		} else {
			panic("Invalid router upstreaming passthrough")
		}
	}
}
