package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"strconv"
)

type peerListener interface {
	Port() int
}

type ServiceMeta map[string]string

// 将Acceptor注册到服务发现,IP自动取本地IP
func Register(p cellnet.Peer, options ...interface{}) *discovery.ServiceDesc {
	host := util.GetLocalIP()

	property := p.(cellnet.PeerProperty)

	sd := &discovery.ServiceDesc{
		ID:   MakeLocalSvcID(property.Name()),
		Name: property.Name(),
		Host: host,
		Port: p.(peerListener).Port(),
	}

	sd.SetMeta("SvcGroup", GetSvcGroup())
	sd.SetMeta("SvcIndex", strconv.Itoa(GetSvcIndex()))

	for _, opt := range options {

		switch optValue := opt.(type) {
		case ServiceMeta:
			for metaKey, metaValue := range optValue {
				sd.SetMeta(metaKey, metaValue)
			}
		}
	}

	if GetWANIP() != "" {
		sd.SetMeta("WANAddress", util.JoinAddress(GetWANIP(), sd.Port))
	}

	log.SetColor("green").Debugf("service '%s' listen at port: %d", sd.ID, sd.Port)

	// 有同名的要先解除注册，再注册，防止watch不触发
	discovery.Default.Deregister(sd.ID)

	p.(cellnet.ContextSet).SetContext("sd", sd)

	err := discovery.Default.Register(sd)
	if err != nil {
		log.Errorf("service register failed, %s %s", sd.String(), err.Error())
	}

	return sd
}

// 解除peer注册
func Unregister(p cellnet.Peer) {
	property := p.(cellnet.PeerProperty)
	discovery.Default.Deregister(MakeLocalSvcID(property.Name()))
}

// 根据进程的互联发现规则和给定的服务名过滤发现的服务
func DiscoveryService(rules []MatchRule, svcName string, matchSvcGroup string) (ret []*discovery.ServiceDesc) {
	descList := discovery.Default.Query(svcName)

	if len(descList) > 0 {

		// 保持服务发现中的所有连接
		for _, sd := range MatchService(rules, svcName, descList) {

			if matchSvcGroup == "" || (matchSvcGroup != "" && sd.GetMeta("SvcGroup") == matchSvcGroup) {
				ret = append(ret, sd)
			}

		}
	}

	return
}

type DiscoveryOption struct {
	MaxCount       int  // 连接数，默认发起多条连接
	ForceSelfGroup bool // 默认只找与自己同组的服务
}

// 发现一个服务，服务可能拥有多个地址，每个地址返回时，创建一个connector并开启
func DiscoveryConnector(rules []MatchRule, tgtSvcName string, opt DiscoveryOption, peerCreator func(*discovery.ServiceDesc) cellnet.Peer) {

	// 从发现到连接有一个过程，需要用Map防止还没连上，又创建一个新的连接
	connectorBySvcID := newConnSet()

	notify := discovery.Default.RegisterNotify("add")
	for {

		var matchSvcGroup string
		if opt.ForceSelfGroup {
			matchSvcGroup = GetSvcGroup()
		}

		descList := DiscoveryService(rules, tgtSvcName, matchSvcGroup)

		// 保持服务发现中的所有连接
		for _, sd := range descList {

			// 同一个svcid，永久对应一个connector，断开后自动重连
			if connectorBySvcID.Get(sd.ID) == nil {

				// 达到最大连接
				if opt.MaxCount > 0 && connectorBySvcID.Count() >= opt.MaxCount {
					continue
				}

				p := peerCreator(sd)

				if p != nil {
					contextSet := p.(cellnet.ContextSet)
					contextSet.SetContext("sd", sd)
					contextSet.SetContext("connSet", connectorBySvcID)
					connectorBySvcID.Add(sd.ID, p)
				}

			}
		}

		<-notify
	}
}
