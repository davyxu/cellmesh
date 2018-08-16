package main

import (
	_ "github.com/davyxu/cellmesh/demo/login/login"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service/cellsvc"
	"github.com/davyxu/cellmesh/svcfx"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	svcfx.Init()

	s := cellsvc.NewService("demo.login")
	s.SetDispatcher(proto.GetDispatcher("demo.login"))

	s.Start()

	util.WaitExit()

	s.Stop()
}
