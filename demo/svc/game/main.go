package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
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

		// 开始接收game状态
		//hubstatus.StartSendStatus("game_status", time.Second*3, func() int {
		//	return 100
		//})
	})

	basefx.StartLoop(nil)

	basefx.Exit()
}
