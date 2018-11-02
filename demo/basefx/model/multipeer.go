package fxmodel

import (
	"fmt"
	"github.com/davyxu/cellnet"
)

type MultiStatus interface {
	GetPeers() []cellnet.Peer
	String() string
}

type MultiPeer struct {
	peers []cellnet.Peer
	param ServiceParameter
}

func (self *MultiPeer) String() string {
	return fmt.Sprintf("%13s %15s", self.param.SvcName, self.param.NetPeerType)
}

func (self *MultiPeer) Start() cellnet.Peer {
	return self
}

func (self *MultiPeer) Stop() {

}

func (self *MultiPeer) TypeName() string {
	return ""
}

func (self *MultiPeer) GetPeers() []cellnet.Peer {
	return self.peers
}

func (self *MultiPeer) IsReady() bool {

	if len(self.peers) == 0 {
		return false
	}

	for _, p := range self.peers {
		if !p.(readyChecker).IsReady() {
			return false
		}
	}

	return true
}

func (self *MultiPeer) Add(p cellnet.Peer) {
	self.peers = append(self.peers, p)
}

func NewMultiPeer(param ServiceParameter) *MultiPeer {
	return &MultiPeer{param: param}
}
