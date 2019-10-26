package link

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
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
		ID:   cellmesh.MakeSvcID(property.Name()),
		Name: property.Name(),
		Host: host,
		Port: p.(peerListener).Port(),
	}

	for _, opt := range options {

		switch optValue := opt.(type) {
		case ServiceMeta:
			for metaKey, metaValue := range optValue {
				sd.SetMeta(metaKey, metaValue)
			}
		}
	}

	// TODO 自动获取外网IP
	if cellmesh.GetWANIP() != "" {
		sd.SetMeta(cellmesh.SDMetaKey_WANAddress, util.JoinAddress(cellmesh.GetWANIP(), sd.Port))
	}

	log.SetColor("green").Debugf("service '%s' listen at port: %d", sd.ID, sd.Port)

	p.(cellnet.ContextSet).SetContext(cellmesh.PeerContextKey_ServiceDesc, sd)

	// 有同名的要先解除注册，再注册，防止watch不触发
	discovery.Default.Deregister(sd.ID)
	err := discovery.Default.Register(sd)
	if err != nil {
		log.Errorf("service register failed, %s %s", sd.String(), err.Error())
	}

	return sd
}

func GetPeerDesc(p cellnet.Peer) (ret *discovery.ServiceDesc) {

	if cs, ok := p.(cellnet.ContextSet); ok {
		if cs.FetchContext(cellmesh.PeerContextKey_ServiceDesc, &ret) {
			return
		}
	}

	return nil
}

// 解除peer注册
func Unregister(p cellnet.Peer) {
	property := p.(cellnet.PeerProperty)
	discovery.Default.Deregister(cellmesh.MakeSvcID(property.Name()))
}
