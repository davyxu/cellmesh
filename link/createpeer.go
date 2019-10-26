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
func StartService(param *ServiceParameter) cellnet.Peer {

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

func addRemoteSvc(desc *discovery.ServiceDesc, param *ServiceParameter) {
	preses := GetLink(desc.ID)
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
	MarkLink(p.(cellnet.TCPConnector).Session(), desc.ID, desc.Name)

	p.Start()

	p.(cellnet.ContextSet).SetContext(cellmesh.PeerContextKey_ServiceDesc, desc)

	// 注册本地Peer，方便做检查遍历
	AddLocalPeer(p)
}

var (
	DebugSyncService bool
)

// 让服务发现和本地peer同步
func syncService(param *ServiceParameter) {

	// 本类服务的所有服务id
	descList := discovery.Default.Query(param.SvcName)
	var peerToRemove []cellnet.Peer

	// 遍历所有的本地服务
	VisitLocalPeer(func(p cellnet.Peer) bool {

		sd := GetPeerDesc(p)
		if DebugSyncService {
			log.Debugf("peer: %+v", sd)
		}

		if sd != nil {

			// 是否在服务发现列表中存在
			if !discovery.DescExistsByID(sd.ID, descList) {
				peerToRemove = append(peerToRemove, p)
			}
		}

		return true
	})

	// 已经没有的服务, 删除
	for _, p := range peerToRemove {
		sd := GetPeerDesc(p)

		log.SetColor("yellow").Infof("remove remote service : '%s'", sd.ID)

		// 有自动连接情况时, 关闭
		p.Stop()
		RemoveLocalPeer(p)
	}

	// 加入已有的服务
	for _, desc := range descList {

		if DebugSyncService {
			log.Debugf("check remote link: %s", desc.String())
		}
		addRemoteSvc(desc, param)
	}

	if DebugSyncService {
		log.Debugf("sync done")
	}
}

// 连接到服务
func LinkService(param *ServiceParameter) {

	// 填充默认值
	param.MakeConnectorDefault()

	// 提前注册回调, 避免在处理已有服务时掉服务
	notify := discovery.Default.RegisterNotify()

	syncService(param)

	// 实时检查更新
	go func() {
		for {

			notifyCtx := <-notify

			if DebugSyncService {
				log.Debugf("recv notify %+v", notifyCtx)
			}

			syncService(param)
		}
	}()

}
