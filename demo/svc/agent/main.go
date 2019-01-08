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
	"github.com/davyxu/cellmesh/demo/svc/hub/api"
	"github.com/davyxu/cellmesh/demo/svc/hub/status"
	"github.com/davyxu/cellmesh/service"
	_ "github.com/davyxu/cellnet/peer/gorillaws"
	"github.com/davyxu/golog"
	"time"
)

var log = golog.New("main")

func main() {

	basefx.Init("agent")

	routerule.Download()

	heartbeat.StartCheck()

	model.AgentSvcID = service.GetLocalSvcID()

	// 要连接的服务列表
	basefx.CreateCommnicateConnector(fxmodel.ServiceParameter{
		SvcName:      "game",
		MaxConnCount: -1,
		NetProcName:  "agent.backend",
	})

	switch *fxmodel.FlagCommunicateType {
	case "tcp":
		frontend.Start(model.FrontendParameter{
			SvcName:     "agent",
			ListenAddr:  ":0",
			NetPeerType: "tcp.Acceptor",
			NetProcName: "tcp.frontend",
		})
	case "ws":
		frontend.Start(model.FrontendParameter{
			SvcName:     "agent",
			ListenAddr:  ":0",
			NetPeerType: "gorillaws.Acceptor",
			NetProcName: "ws.frontend",
		})
	}

	hubapi.ConnectToHub(func() {

		// 发送网关连接数量
		hubstatus.StartSendStatus("agent_status", time.Second*3, func() int {
			return model.FrontendSessionManager.SessionCount()
		})
	})

	basefx.StartLoop(nil)

	frontend.Stop()

	basefx.Exit()
}
