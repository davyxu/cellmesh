package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
	_ "github.com/davyxu/cellmesh/demo/login/login"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	basefx.Init("login")

	basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
		SvcName:    "login",
		ListenAddr: ":0",
	})

	basefx.StartLoop()

	basefx.Exit()
}
