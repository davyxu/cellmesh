package main

import (
	_ "github.com/davyxu/cellmesh/demo/agent/backend"
	"github.com/davyxu/cellmesh/demo/agent/frontend"
	"github.com/davyxu/cellmesh/demo/agent/routerule"
	"github.com/davyxu/cellmesh/demo/proto"
	_ "github.com/davyxu/cellmesh/demo/proto" // 进入协议
	"github.com/davyxu/cellmesh/service/cellsvc"
	"github.com/davyxu/cellmesh/svcfx"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	svcfx.Init()

	routerule.Download()

	s := cellsvc.NewService("router")
	s.SetDispatcher(proto.GetDispatcher("router"))
	s.Start()

	frontend.Start(":18000")

	util.WaitExit()

	frontend.Stop()
	s.Stop()
}
