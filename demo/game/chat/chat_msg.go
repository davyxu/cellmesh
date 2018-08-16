package chat

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
)

func init() {

	proto.Handle_Game_ChatREQ = func(ev service.Event, req *proto.ChatREQ) {
		fmt.Printf("chat: %+v \n", req.Content)
		ev.Reply(&proto.ChatACK{
			Content: req.Content,
		})
	}
}
