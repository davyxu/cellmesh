package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/hub/api"
	"github.com/davyxu/cellmesh/demo/svc/hub/status"
	_ "github.com/davyxu/cellmesh/demo/svc/login/login"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	basefx.Init("login")

	basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
		SvcName:    "login",
		ListenAddr: ":0",
	})

	hubapi.ConnectToHub(func() {

		// 开始接收game状态
		hubstatus.StartRecvStatus("game_status", &proto.Handle_Login_SvcStatusACK)
	})

	basefx.StartLoop(nil)

	basefx.Exit()
}
