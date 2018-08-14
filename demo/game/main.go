package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/game/chat"
	"github.com/davyxu/cellmesh/demo/game/verify"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
)

func main() {

	svcfx.Init()

	//go service.PrepareConnection("demo.router", service.NewRPCRequestor, nil)

	s := service.NewService("demo.game")

	proto.Register_Verify(s, verify.Verify)
	proto.Register_Chat(s, chat.Chat)

	err := s.Run()

	if err != nil {
		fmt.Println(err)
	}
}
