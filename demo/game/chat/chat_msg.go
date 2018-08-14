package chat

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
)

func Chat(req *proto.ChatREQ, ack *proto.ChatACK) {

	fmt.Printf("chat: %+v \n", req.Content)

	ack.Content = req.Content
}
