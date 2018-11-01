package chat

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/agent/api"
	"github.com/davyxu/cellnet"
)

func init() {

	proto.Handle_Game_ChatREQ = agentapi.HandleBackendMessage(func(ev cellnet.Event, cid proto.ClientID) {

		msg := ev.Message().(*proto.ChatREQ)

		fmt.Printf("chat: %+v \n", msg.Content)

		agentapi.BroadcastAll(&proto.ChatACK{
			Content: msg.Content,
		})
	})
}
