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

		// 消息广播到网关并发给客户端
		agentapi.BroadcastAll(&proto.ChatACK{
			Content: msg.Content,
		})

		// 消息单发给客户端
		agentapi.Send(&cid, &proto.TestACK{
			Dummy: "single send",
		})
	})
}
