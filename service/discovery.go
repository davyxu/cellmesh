package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
)

type DiscoveryOption struct {
	Rules         []MatchRule
	MaxCount      int    // 连接数，默认发起多条连接
	MatchSvcGroup string // 空时，匹配所有同类服务，否则找指定组的服务
}

// 发现一个服务，服务可能拥有多个地址，每个地址返回时，创建一个connector并开启
// DiscoveryService返回值返回持有多个Peer的peer, 判断Peer的IsReady可以得到所有连接准备好的状态
func DiscoveryService(tgtSvcName string, opt DiscoveryOption, peerCreator func(*discovery.ServiceDesc) cellnet.Peer) cellnet.Peer {

	// 从发现到连接有一个过程，需要用Map防止还没连上，又创建一个新的连接
	multiPeer := newMultiPeer()

	go func() {

		notify := discovery.Default.RegisterNotify("add")
		for {

			QueryService(tgtSvcName,
				Filter_MatchRule(opt.Rules),
				Filter_MatchSvcGroup(opt.MatchSvcGroup),
				func(desc *discovery.ServiceDesc) interface{} {

					prePeer := multiPeer.GetPeer(desc.ID)

					// 如果svcid重复汇报, 可能svcid内容有变化
					if prePeer != nil {

						var preDesc *discovery.ServiceDesc
						if prePeer.(cellnet.ContextSet).FetchContext("sd", &preDesc) && !preDesc.Equals(desc) {

							log.Infof("service '%s' change desc, %+v -> %+v...", desc.ID, preDesc, desc)

							// 移除之前的连接
							multiPeer.RemovePeer(desc.ID)

							// 停止重连
							prePeer.Stop()

						} else {
							return true
						}

					}

					// 达到最大连接
					if opt.MaxCount > 0 && len(multiPeer.GetPeers()) >= opt.MaxCount {
						return true
					}

					// 用户创建peer
					p := peerCreator(desc)

					if p != nil {
						contextSet := p.(cellnet.ContextSet)
						contextSet.SetContext("sd", desc)
						multiPeer.AddPeer(desc.ID, p)
					}

					return true
				})

			<-notify
		}

	}()

	return multiPeer
}
