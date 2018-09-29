package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	_ "github.com/davyxu/cellmesh/demo/game/chat"
	_ "github.com/davyxu/cellmesh/demo/game/verify"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/kvconfig"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	service.Init("game")

	basefx.CreateCommnicateAcceptor("game", kvconfig.String(discovery.Default, "config/addr_game", ":10100~10199"))

	util.WaitExit()

	basefx.StopAllPeers()
}
