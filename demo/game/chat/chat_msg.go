package chat

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
)

func Chat(ev *service.Event, req *proto.ChatREQ, ack *proto.ChatACK) {

	fmt.Printf("chat: %+v \n", req.Content)

	ack.Content = req.Content
}
