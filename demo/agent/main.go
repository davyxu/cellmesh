package main

import (
	_ "github.com/davyxu/cellmesh/demo/agent/backend"
	"github.com/davyxu/cellmesh/demo/agent/frontend"
	"github.com/davyxu/cellmesh/demo/agent/heartbeat"
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/demo/agent/routerule"
	"github.com/davyxu/cellmesh/demo/basefx"
	_ "github.com/davyxu/cellmesh/demo/proto" // 进入协议
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/kvconfig"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	service.Init("agent")

	routerule.Download()

	heartbeat.StartCheck()

	basefx.CreateCommnicateAcceptor(model.BackendName, ":0")

	frontend.Start(kvconfig.String(discovery.Default, "cm_demo/config/agent/frontend_addr", ":8001~8101"))

	util.WaitExit()

	frontend.Stop()
	basefx.StopAllPeers()
}
