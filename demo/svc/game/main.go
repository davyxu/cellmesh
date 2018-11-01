package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/demo/proto"
	_ "github.com/davyxu/cellmesh/demo/svc/game/chat"
	_ "github.com/davyxu/cellmesh/demo/svc/game/verify"
	"github.com/davyxu/cellmesh/demo/svc/hub/api"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/kvconfig"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	basefx.Init("game")

	basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
		SvcName:    "game",
		ListenAddr: kvconfig.String(discovery.Default, "config_demo/addr_game", ":0"),
	})

	hubapi.ConnectToHub(func() {
		hubapi.Publish("game_status", &proto.SvcStatusACK{
			UserCount: 110,
		})
	})

	basefx.StartLoop()

	basefx.Exit()
}
