package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
	_ "github.com/davyxu/cellmesh/demo/svc/hub/subscribe"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/kvconfig"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	basefx.Init("hub")

	basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
		SvcName:     "hub",
		NetProcName: "tcp.svc",
		ListenAddr:  kvconfig.String(discovery.Default, "config_demo/addr_hub", ":0"),
	})

	basefx.StartLoop(nil)

	basefx.Exit()
}
