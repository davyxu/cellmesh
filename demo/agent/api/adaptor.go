package api

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
)

// 传入用户处理网关消息回调,返回消息源回调
func HandleRouteMessage(userHandler func(ev service.Event, cid proto.ClientID)) func(ev service.Event) {

	return func(ev service.Event) {

		if cid, ok := ev.PassThrough().(*proto.ClientID); ok {
			userHandler(ev, *cid)
		} else {
			panic("Invalid router upstreaming passthrough")
		}
	}
}
