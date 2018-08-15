package main

import (
	"github.com/davyxu/cellmesh/demo/game/chat"
	"github.com/davyxu/cellmesh/demo/game/verify"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/service/cellsvc"
	"github.com/davyxu/cellmesh/svcfx"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	svcfx.Init()

	dis := service.NewDispatcher()
	proto.Serve_Verify(dis, verify.Verify)
	proto.Serve_Chat(dis, chat.Chat)

	s := cellsvc.NewService("demo.game")
	s.SetDispatcher(dis)
	s.Start()
	sd := s.(interface {
		GetSD() *discovery.ServiceDesc
	}).GetSD()

	r := cellsvc.NewConnService("demo.router", sd)
	r.SetDispatcher(dis)
	r.Start()

	util.WaitExit()

	r.Stop()
	s.Stop()
}
