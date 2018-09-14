package main

import (
	_ "github.com/davyxu/cellmesh/demo/game/chat"
	_ "github.com/davyxu/cellmesh/demo/game/verify"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service/cellsvc"
	"github.com/davyxu/cellmesh/svcfx"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	svcfx.Init()

	con := cellsvc.NewConnector("svc->agent")
	con.SetProcessor("tcp.ltv")
	con.SetEventCallback(proto.GetDispatcher("game"))
	con.Start()

	util.WaitExit()

	con.Stop()
}
