package main

import (
	"github.com/davyxu/cellmesh/demo/agent/backend"
	"github.com/davyxu/cellmesh/demo/agent/frontend"
	"github.com/davyxu/cellmesh/demo/proto"
	_ "github.com/davyxu/cellmesh/demo/proto" // 进入协议
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

	s := cellsvc.NewService("demo.router")
	s.SetDispatcher(dis)
	proto.Serve_RouterBindUser(dis, backend.RouterBindUser)
	s.Start()

	frontend.Start()

	util.WaitExit()

	frontend.Stop()
	s.Stop()
}
