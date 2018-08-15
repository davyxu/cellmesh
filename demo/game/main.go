package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/game/chat"
	"github.com/davyxu/cellmesh/demo/game/verify"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	svcfx.Init()

	s := service.NewService("demo.game")

	proto.Serve_Verify(s, verify.Verify)
	proto.Serve_Chat(s, chat.Chat)

	err := s.Run()

	if err != nil {
		fmt.Println(err)
	}
}
