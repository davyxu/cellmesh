package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/hub/api"
	_ "github.com/davyxu/cellmesh/demo/svc/login/login"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	basefx.Init("login")

	basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
		SvcName:    "login",
		ListenAddr: ":0",
	})

	proto.Handle_Login_SvcStatusACK = func(ev cellnet.Event) {

		msg := ev.Message().(*proto.SvcStatusACK)
		log.Debugln(msg.UserCount)
	}

	hubapi.ConnectToHub(func() {
		hubapi.Subscribe("game_status")
	})

	basefx.StartLoop()

	basefx.Exit()
}
