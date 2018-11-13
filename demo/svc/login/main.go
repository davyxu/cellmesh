package main

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/hub/api"
	"github.com/davyxu/cellmesh/demo/svc/hub/status"
	_ "github.com/davyxu/cellmesh/demo/svc/login/login"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/peer/gorillaws"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/gorillaws"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	basefx.Init("login")

	// 与客户端通信的处理器
	proc.RegisterProcessor("ws.client", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(new(gorillaws.WSMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(gorillaws.MsgHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})

	switch *fxmodel.FlagCommunicateType {
	case "tcp":
		basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
			SvcName:     "login",
			NetPeerType: "tcp.Acceptor",
			NetProcName: "tcp.client",
			ListenAddr:  ":0",
		})
	case "ws":
		basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
			SvcName:     "login",
			NetPeerType: "gorillaws.Acceptor",
			NetProcName: "ws.client",
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
