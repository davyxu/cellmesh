package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/hub/api"
	"github.com/davyxu/cellmesh/demo/svc/hub/status"
	_ "github.com/davyxu/cellmesh/demo/svc/login/login"
	_ "github.com/davyxu/cellnet/peer/gorillaws"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	basefx.Init("login")

	switch *fxmodel.FlagCommunicateType {
	case "tcp":
		basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
			SvcName:     "login",
			NetPeerType: "tcp.Acceptor",
			NetProcName: "tcp.svc",
			ListenAddr:  ":0",
		})
	case "ws":
		basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
			SvcName:     "login",
			NetPeerType: "gorillaws.Acceptor",
			NetProcName: "ws.svc",
			ListenAddr:  ":0",
		})
	}

	hubapi.ConnectToHub(func() {

		// 开始接收game状态
		hubstatus.StartRecvStatus([]string{"game_status", "agent_status"}, &proto.Handle_Login_SvcStatusACK)
	})

	basefx.StartLoop(nil)

	basefx.Exit()
}
