package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
	_ "github.com/davyxu/cellmesh/demo/proto" // 进入协议
	_ "github.com/davyxu/cellmesh/demo/svc/agent/backend"
	"github.com/davyxu/cellmesh/demo/svc/agent/frontend"
	"github.com/davyxu/cellmesh/demo/svc/agent/heartbeat"
	"github.com/davyxu/cellmesh/demo/svc/agent/model"
	"github.com/davyxu/cellmesh/demo/svc/agent/routerule"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/kvconfig"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	basefx.Init("agent")

	routerule.Download()

	heartbeat.StartCheck()

	model.AgentSvcID = service.MakeLocalSvcID(model.BackendName)

	// 要连接的服务列表
	basefx.CreateCommnicateConnector(fxmodel.ServiceParameter{
		SvcName:      "game",
		MaxConnCount: -1,
	})

	frontend.Start(kvconfig.String(discovery.Default, "config_demo/addr_agentfrontend", ":8001~8101"))

	basefx.StartLoop()

	frontend.Stop()

	basefx.Exit()
}
