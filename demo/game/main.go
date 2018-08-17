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

	r := cellsvc.NewConnector("game", "router")
	r.SetDispatcher(proto.GetDispatcher("game"))
	r.Start()

	util.WaitExit()

	r.Stop()
}
