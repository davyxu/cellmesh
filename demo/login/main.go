package main

import (
	"github.com/davyxu/cellmesh/demo/login/login"
	"github.com/davyxu/cellmesh/demo/proto"
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

	s := cellsvc.NewService("demo.login")
	s.SetDispatcher(dis)

	proto.Serve_Login(dis, login.Login)

	s.Start()

	util.WaitExit()

	s.Stop()
}
