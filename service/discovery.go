package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"sync"
)

type peerListener interface {
	Port() int
}

// 将Acceptor注册到服务发现,IP自动取本地IP
func Register(p cellnet.Peer) {
	host := util.GetLocalIP()

	property := p.(cellnet.PeerProperty)

	sd := &discovery.ServiceDesc{
		Host: host,
		Port: p.(peerListener).Port(),
		ID:   MakeLocalSvcID(property.Name()),
		Name: property.Name(),
		Tags: []string{GetSvcGroup()},
	}

	if GetWANIP() != "" {
		sd.SetMeta("WANAddress", util.JoinAddress(GetWANIP(), sd.Port))
	}

	log.SetColor("green").Debugf("service '%s' listen at %s:%d", sd.ID, host, sd.Port)

	// 有同名的要先解除注册，再注册，防止watch不触发
	discovery.Default.Deregister(sd.ID)

	discovery.Default.Register(sd)
}

// 解除peer注册
func Unregister(p cellnet.Peer) {
	property := p.(cellnet.PeerProperty)
	discovery.Default.Deregister(MakeLocalSvcID(property.Name()))
}

// 根据进程的互联发现规则和给定的服务名过滤发现的服务
func DiscoveryService(rules []MatchRule, svcName string) (ret []*discovery.ServiceDesc) {
	descList, err := discovery.Default.Query(svcName)

	if err == nil && len(descList) > 0 {

		// 保持服务发现中的所有连接
		for _, sd := range MatchService(rules, svcName, descList) {
			ret = append(ret, sd)
		}
	}

	return
}

// 发现一个服务，服务可能拥有多个地址，每个地址返回时，创建一个connector并开启
func DiscoveryConnector(rules []MatchRule, tgtSvcName string, peerCreator func(*discovery.ServiceDesc) cellnet.Peer) {

	// 从发现到连接有一个过程，需要用Map防止还没连上，又创建一个新的连接
	var connectorBySvcID sync.Map

	notify := discovery.Default.RegisterNotify("add")
	for {

		descList := DiscoveryService(rules, tgtSvcName)

		// 保持服务发现中的所有连接
		for _, sd := range descList {

			// 新连接马上连接，老连接保留
			if _, ok := connectorBySvcID.Load(sd.ID); !ok {

				p := peerCreator(sd)

				contextSet := p.(cellnet.ContextSet)
				contextSet.SetContext("sd", sd)
				contextSet.SetContext("connMap", &connectorBySvcID)

				connectorBySvcID.Store(sd.ID, p)
			}
		}

		<-notify
	}
}
