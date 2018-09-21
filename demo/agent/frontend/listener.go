package frontend

import (
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"time"
)

func Start(addr string) {

	clientListener := peer.NewGenericPeer("tcp.Acceptor", model.FrontendName, addr, nil)

	proc.BindProcessorHandler(clientListener, "tcp.frontend", nil)

	socketOpt := clientListener.(cellnet.TCPSocketOption)

	// 无延迟设置缓冲
	socketOpt.SetSocketBuffer(2048, 2048, true)

	// 40秒无读，20秒无写断开
	socketOpt.SetSocketDeadline(time.Second*40, time.Second*20)

	clientListener.Start()
	model.FrontendSessionManager = clientListener.(peer.SessionManager)

	model.AgentSvcID = service.MakeLocalSvcID(model.FrontendName)

	// 服务发现注册服务
	service.Register(clientListener)
}

func Stop() {

	if model.FrontendSessionManager != nil {
		model.FrontendSessionManager.(cellnet.Peer).Stop()
		discovery.Default.Deregister(model.AgentSvcID)
	}

}
