package linkmgr

import (
	"github.com/davyxu/cellmesh/discovery"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"time"
)

func StartService(param ServiceParameter) cellnet.Peer {

	param.MakeServiceDefault()

	p := peer.NewGenericPeer(param.PeerType, param.SvcName, param.ListenAddress, param.Queue)

	//"tcp.svc"
	proc.BindProcessorHandler(p, param.NetProc, param.EventCallback)

	if opt, ok := p.(cellnet.TCPSocketOption); ok {
		opt.SetSocketBuffer(2048, 2048, true)
	}

	AddLocalPeer(p)

	p.Start()

	Register(p)

	return p
}

func ConnectToService(param ServiceParameter) {

	param.MakeConnectorDefault()

	discovery.Default.SetNotify(func(evType string, args ...interface{}) {

		if evType == "add" {
			discovery.QueryService(param.SvcName, func(desc *discovery.ServiceDesc) interface{} {

				// TODO 全局连接约束表

				preses := GetRemoteSession(desc.ID)
				if preses != nil {

					// TODO 更换地址操作
					return true
				}

				p := peer.NewGenericPeer(param.PeerType, param.SvcName, desc.Address(), param.Queue)

				proc.BindProcessorHandler(p, param.NetProc, param.EventCallback)

				if opt, ok := p.(cellnet.TCPSocketOption); ok {
					opt.SetSocketBuffer(2048, 2048, true)
				}

				p.(cellnet.TCPConnector).SetReconnectDuration(time.Second * 3)

				// 提前添加到表中， 通过IsReady判断，避免连不上时反复创建Connector
				AddRemoteSession(p.(cellnet.TCPConnector).Session(), desc.ID)

				p.Start()

				// 注册本地Peer，方便做检查遍历
				AddLocalPeer(p)

				return true
			})
		}

	})

}
