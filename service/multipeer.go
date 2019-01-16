package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"sync"
)

type MultiPeer interface {
	GetPeers() []cellnet.Peer
	cellnet.ContextSet
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

func (self *multiPeer) AddPeer(svcid string, p cellnet.Peer) {
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
