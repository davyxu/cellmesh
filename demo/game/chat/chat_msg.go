package chat

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/agent/api"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellnet"
)

func init() {

	proto.Handle_Game_ChatREQ = api.HandleBackendMessage(func(ev cellnet.Event, cid proto.ClientID) {

		msg := ev.Message().(*proto.ChatREQ)

		fmt.Printf("chat: %+v \n", msg.Content)

		api.BroadcastAll(&proto.ChatACK{
			Content: msg.Content,
		})
	})
}
