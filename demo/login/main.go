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

	acc := cellsvc.NewAcceptor("login")
	acc.SetProcessor("tcp.ltv")
	acc.SetEventCallback(proto.GetDispatcher("login"))
	acc.Start()

	util.WaitExit()

	acc.Stop()
}
