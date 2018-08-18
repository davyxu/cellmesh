package chat

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/agent/api"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
)

func init() {

	proto.Handle_Game_ChatREQ = api.HandleRouteMessage(func(ev service.Event, cid proto.ClientID) {

		msg := ev.Message().(*proto.ChatREQ)

		fmt.Printf("chat: %+v \n", msg.Content)

		api.BroadcastAll(&proto.ChatACK{
			Content: msg.Content,
		})
	})
}
