package link

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"time"
)

func addRemoteSvc(desc *discovery.ServiceDesc, param *ServiceParameter) {

	p := peer.NewGenericPeer(param.PeerType, param.SvcName, desc.Address(), param.Queue)

	proc.BindProcessorHandler(p, param.NetProc, param.EventCallback)

	if opt, ok := p.(cellnet.TCPSocketOption); ok {
		opt.SetSocketBuffer(2048, 2048, true)
	}

	p.(cellnet.TCPConnector).SetReconnectDuration(time.Second * 3)

	p.(cellnet.ContextSet).SetContext(cellmesh.PeerContextKey_ServiceDesc, desc)

	// 提前添加到表中， 通过IsReady判断，避免连不上时反复创建Connector
	markLink(p.(cellnet.TCPConnector).Session(), desc.ID, desc.Name)

	p.Start()

	// 注册本地Peer，方便做检查遍历
	AddPeer(p)
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
	VisitPeer(func(p cellnet.Peer) bool {

		sd := GetPeerDesc(p)
		if DebugSyncService {
			log.Debugf("peer: %+v", sd)
		}

		if sd != nil {

			// 只处理同类的服务
			if sd.Name != param.SvcName {
				return true
			}

			// 是否在服务发现列表中存在
			if !discovery.DescExistsByID(sd.ID, descList) {
				peerToRemove = append(peerToRemove, p)
			}
		}

		return true
	})

	// 已经没有的服务, 删除
	for _, p := range peerToRemove {
		// 有自动连接情况时, 关闭
		p.Stop()
		RemovePeer(p)
	}

	// 加入已有的服务
	for _, desc := range descList {

		if DebugSyncService {
			log.Debugf("check remote link: %s", desc.String())
		}

		preses := GetLink(desc.ID)
		if preses == nil {

			dummy := GetPeer(param.SvcName)
			if _, ok := dummy.(interface {
				IsDummy() bool
			}); ok {
				RemovePeer(dummy)
			}

			addRemoteSvc(desc, param)
		}
	}

	if DebugSyncService {
		log.Debugf("sync done")
	}
}

// 连接到服务
func ConnectService(param *ServiceParameter) {
	AddPeer(newDummyPeer(param))

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
