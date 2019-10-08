package link

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"time"
)

// 开启服务
func StartService(param ServiceParameter) cellnet.Peer {

	// 填充默认值
	param.MakeServiceDefault()

	p := peer.NewGenericPeer(param.PeerType, param.SvcName, param.ListenAddress, param.Queue)

	proc.BindProcessorHandler(p, param.NetProc, param.EventCallback)

	if opt, ok := p.(cellnet.TCPSocketOption); ok {
		opt.SetSocketBuffer(2048, 2048, true)
	}

	AddLocalPeer(p)

	p.Start()

	Register(p)

	return p
}

func addRemoteSvc(desc *discovery.ServiceDesc, param ServiceParameter) {
	if GetLocalPeer(desc.Name) != nil {
		return
	}

	log.Debugf("found svc add '%s'", desc.String())

	// TODO 全局连接约束表

	preses := GetRemoteLink(desc.ID)
	if preses != nil {
		return
	}

	p := peer.NewGenericPeer(param.PeerType, param.SvcName, desc.Address(), param.Queue)

	proc.BindProcessorHandler(p, param.NetProc, param.EventCallback)

	if opt, ok := p.(cellnet.TCPSocketOption); ok {
		opt.SetSocketBuffer(2048, 2048, true)
	}

	p.(cellnet.TCPConnector).SetReconnectDuration(time.Second * 3)

	// 提前添加到表中， 通过IsReady判断，避免连不上时反复创建Connector
	AddRemoteLink(p.(cellnet.TCPConnector).Session(), desc.ID, desc.Name)

	p.Start()

	p.(cellnet.ContextSet).SetContext(cellmesh.PeerContextKey_ServiceDesc, desc)

	// 注册本地Peer，方便做检查遍历
	AddLocalPeer(p)
}

// 连接到服务
func LinkService(param ServiceParameter) {

	// 填充默认值
	param.MakeConnectorDefault()

	// 提前注册回调, 避免在处理已有服务时掉服务
	notify := discovery.Default.RegisterNotify()

	// 加入已有的服务
	for _, desc := range discovery.Default.Query(param.SvcName) {
		addRemoteSvc(desc, param)
	}

	// 实时检查更新
	go func() {
		for {

			notifyDesc := <-notify

			desc := notifyDesc.Desc

			// 根据服务变更方式响应
			switch notifyDesc.Mode {
			case "add": // 新增服务
				addRemoteSvc(desc, param)

			case "del": // 移除服务

				log.Debugf("found svc del '%s'", desc.ID)

				peerToRemove := getLocalPeerBySvcID(desc.ID)

				if peerToRemove != nil {
					peerToRemove.Stop()
					RemoveLocalPeer(peerToRemove)
				}

			case "mod": // 服务信息修改 TODO 更换地址操作
			}

		}
	}()

}

func getLocalPeerBySvcID(svcID string) (ret cellnet.Peer) {

	VisitLocalPeer(func(i cellnet.Peer) bool {
		var desc *discovery.ServiceDesc
		if i.(cellnet.ContextSet).FetchContext(cellmesh.PeerContextKey_ServiceDesc, &desc) {
			if desc.ID == svcID {
				ret = i
				return false
			}
		}

		return true
	})

	return
}
