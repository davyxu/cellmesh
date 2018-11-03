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
	} else {
		// 没有设外网IP时, 用本地IP
		sd.SetMeta("WANAddress", util.JoinAddress(util.GetLocalIP(), sd.Port))
	}

	log.SetColor("green").Debugf("service '%s' listen at port: %d", sd.ID, sd.Port)

	p.(cellnet.ContextSet).SetContext("sd", sd)

	// 有同名的要先解除注册，再注册，防止watch不触发
	discovery.Default.Deregister(sd.ID)
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
