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
func DiscoveryConnector(tgtSvcName string, opt DiscoveryOption, peerCreator func(*discovery.ServiceDesc) cellnet.Peer) {

	// 从发现到连接有一个过程，需要用Map防止还没连上，又创建一个新的连接
	connectorBySvcID := newConnSet()

	notify := discovery.Default.RegisterNotify("add")
	for {

		QueryService(tgtSvcName,
			Filter_MatchRule(opt.Rules),
			Filter_MatchSvcGroup(opt.MatchSvcGroup),
			func(desc *discovery.ServiceDesc) interface{} {

				// 同一个svcid，永久对应一个connector，断开后自动重连
				if connectorBySvcID.Get(desc.ID) == nil {

					// 达到最大连接
					if opt.MaxCount > 0 && connectorBySvcID.Count() >= opt.MaxCount {
						return true
					}

					p := peerCreator(desc)

					if p != nil {
						contextSet := p.(cellnet.ContextSet)
						contextSet.SetContext("sd", desc)
						contextSet.SetContext("connSet", connectorBySvcID)
						connectorBySvcID.Add(desc.ID, p)
					}

				}

				return true
			})

		<-notify
	}
}
