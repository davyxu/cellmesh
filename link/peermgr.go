package link

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"sync"
)

var (
	linkPeers      []cellnet.Peer
	linkPeersGuard sync.RWMutex
)

type linkPeerProp interface {
	Name() string
	Session() cellnet.Session
}

func AddPeer(p cellnet.Peer) {
	linkPeersGuard.Lock()
	linkPeers = append(linkPeers, p)
	linkPeersGuard.Unlock()
}

func GetPeerLink(peer cellnet.Peer) cellnet.Session {
	if lp, ok := peer.(linkPeerProp); ok {
		return lp.Session()
	}

	return nil
}

func RemovePeer(p cellnet.Peer) {
	linkPeersGuard.Lock()
	for index, libp := range linkPeers {
		if libp == p {

			linkPeers = append(linkPeers[:index], linkPeers[index+1:]...)
			break
		}

	}
	linkPeersGuard.Unlock()
}

func GetPeer(svcName string) cellnet.Peer {

	linkPeersGuard.RLock()
	defer linkPeersGuard.RUnlock()

	for _, svc := range linkPeers {
		if prop, ok := svc.(linkPeerProp); ok && prop.Name() == svcName {
			return svc
		}
	}

	return nil
}

func VisitPeer(callback func(cellnet.Peer) bool) {
	linkPeersGuard.RLock()
	defer linkPeersGuard.RUnlock()

	for _, svc := range linkPeers {
		if !callback(svc) {
			break
		}
	}
}

func StopPeers() {
	linkPeersGuard.RLock()
	defer linkPeersGuard.RUnlock()

	for i := len(linkPeers) - 1; i >= 0; i-- {
		svc := linkPeers[i]
		svc.Stop()
	}
}

func GetPeerDesc(p cellnet.Peer) (ret *discovery.ServiceDesc) {

	if cs, ok := p.(cellnet.ContextSet); ok {
		if cs.FetchContext(cellmesh.PeerContextKey_ServiceDesc, &ret) {
			return
		}
	}

	return nil
}
