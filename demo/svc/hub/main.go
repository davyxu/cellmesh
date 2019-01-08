package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
	_ "github.com/davyxu/cellmesh/demo/svc/hub/subscribe"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	basefx.Init("hub")

	basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
		SvcName:     "hub",
		NetProcName: "tcp.svc",
		ListenAddr:  ":0",
	})

	basefx.StartLoop(nil)

	basefx.Exit()
}
