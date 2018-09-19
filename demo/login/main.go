package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	_ "github.com/davyxu/cellmesh/demo/login/login"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	service.Init("login")

	basefx.CreateCommnicateAcceptor("login", ":0")

	util.WaitExit()

	basefx.StopAllPeers()
}
