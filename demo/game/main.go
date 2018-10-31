package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
	_ "github.com/davyxu/cellmesh/demo/game/chat"
	_ "github.com/davyxu/cellmesh/demo/game/verify"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/kvconfig"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	basefx.Init("game")

	basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
		SvcName:    "game",
		ListenAddr: kvconfig.String(discovery.Default, "config/addr_game", ":10100~10199"),
	})

	basefx.StartLoop()

	basefx.Exit()
}
