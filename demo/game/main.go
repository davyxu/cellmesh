package main

import (
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/demo/basefx"
	_ "github.com/davyxu/cellmesh/demo/game/chat"
	_ "github.com/davyxu/cellmesh/demo/game/verify"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	service.Init("game")

	basefx.CreateCommnicateConnector(model.BackendName)

	util.WaitExit()

	basefx.StopAllPeers()
}
