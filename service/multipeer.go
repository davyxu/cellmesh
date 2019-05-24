package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"sync"
)

// 一类服务发起多个连接(不是同一地址), 比如 login1 login2
type MultiPeer interface {
	GetPeers() []cellnet.Peer

	cellnet.ContextSet

	AddPeer(sd *discovery.ServiceDesc, p cellnet.Peer)
}

type multiPeer struct {
	peer.CoreContextSet
	peers      []cellnet.Peer
	peersGuard sync.RWMutex
	context    interface{}
}

func (self *multiPeer) Start() cellnet.Peer {
	return self
}

func (self *multiPeer) Stop() {

}

func (self *multiPeer) TypeName() string {
	return ""
}

func (self *multiPeer) GetPeers() []cellnet.Peer {
	self.peersGuard.RLock()
	defer self.peersGuard.RUnlock()

	return self.peers
}

func (self *multiPeer) IsReady() bool {

	peers := self.GetPeers()

	if len(peers) == 0 {
		return false
	}

	for _, p := range peers {
		if !p.(cellnet.PeerReadyChecker).IsReady() {
			return false
		}
	}

	return true
}

// 保证AddPeer在Peer  Start之前调用, 否则在连接上时因为没有sd,会导致不汇报服务信息
func (self *multiPeer) AddPeer(sd *discovery.ServiceDesc, p cellnet.Peer) {

	contextSet := p.(cellnet.ContextSet)
	contextSet.SetContext("sd", sd)

	self.peersGuard.Lock()
	self.peers = append(self.peers, p)
	self.peersGuard.Unlock()
}

func (self *multiPeer) GetPeer(svcid string) cellnet.Peer {
	for _, p := range self.peers {

		if getSvcIDByPeer(p) == svcid {
			return p
		}
	}

	return nil
}

func (self *multiPeer) RemovePeer(svcid string) {
	self.peersGuard.Lock()
	defer self.peersGuard.Unlock()
	for index, p := range self.peers {

		if getSvcIDByPeer(p) == svcid {
			self.peers = append(self.peers[:index], self.peers[index+1:]...)
			break
		}
	}
}

func getSvcIDByPeer(p cellnet.Peer) string {
	var sd *discovery.ServiceDesc
	if p.(cellnet.ContextSet).FetchContext("sd", &sd) {
		return sd.ID
	}

	return ""
}

func newMultiPeer() *multiPeer {
	return &multiPeer{}
}
