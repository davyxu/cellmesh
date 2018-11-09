package frontend

import (
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/demo/svc/agent/model"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"time"
)

func Start(param model.FrontendParameter) {

	clientListener := peer.NewGenericPeer(param.NetPeerType, "agent", param.ListenAddr, nil)

	proc.BindProcessorHandler(clientListener, param.NetProcName, nil)

	if socketOpt, ok := clientListener.(cellnet.TCPSocketOption); ok {
		// 无延迟设置缓冲
		socketOpt.SetSocketBuffer(2048, 2048, true)

		// 40秒无读，20秒无写断开
		socketOpt.SetSocketDeadline(time.Second*40, time.Second*20)
	}

	clientListener.Start()
	model.FrontendSessionManager = clientListener.(peer.SessionManager)

	// 服务发现注册服务
	service.Register(clientListener)

	fxmodel.AddLocalService(clientListener)
}

func Stop() {

	if model.FrontendSessionManager != nil {
		model.FrontendSessionManager.(cellnet.Peer).Stop()
		discovery.Default.Deregister(model.AgentSvcID)
	}

}
