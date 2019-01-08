package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
	_ "github.com/davyxu/cellmesh/demo/svc/game/chat"
	_ "github.com/davyxu/cellmesh/demo/svc/game/verify"
	"github.com/davyxu/cellmesh/demo/svc/hub/api"
	"github.com/davyxu/cellmesh/demo/svc/hub/status"
	"github.com/davyxu/golog"
	"time"
)

var log = golog.New("main")

func main() {

	basefx.Init("game")

	basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
		SvcName:     "game",
		NetProcName: "svc.backend",
		ListenAddr:  ":0",
	})

	hubapi.ConnectToHub(func() {

		// 开始接收game状态
		hubstatus.StartSendStatus("game_status", time.Second*3, func() int {
			return 100
		})
	})

	basefx.StartLoop(nil)

	basefx.Exit()
}
