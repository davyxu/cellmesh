package link

import (
	"github.com/davyxu/cellnet"
	"sync"
)

var (
	localPeers      []cellnet.Peer
	localPeersGuard sync.RWMutex
)

func AddLocalPeer(p cellnet.Peer) {
	localPeersGuard.Lock()
	localPeers = append(localPeers, p)
	localPeersGuard.Unlock()

	log.Debugf("add local peer '%s'", p.(cellnet.PeerProperty).Name())
}

func RemoveLocalPeer(p cellnet.Peer) {
	localPeersGuard.Lock()
	for index, libp := range localPeers {
		if libp == p {
			log.Debugf("remove local peer '%s'", p.(cellnet.PeerProperty).Name())
			localPeers = append(localPeers[:index], localPeers[index+1:]...)
			break
		}

	}
	localPeersGuard.Unlock()
}

func GetLocalPeer(svcName string) cellnet.Peer {

	localPeersGuard.RLock()
	defer localPeersGuard.RUnlock()

	for _, svc := range localPeers {
		if prop, ok := svc.(cellnet.PeerProperty); ok && prop.Name() == svcName {
			return svc
		}
	}

	return nil
}

func VisitLocalPeer(callback func(cellnet.Peer) bool) {
	localPeersGuard.RLock()
	defer localPeersGuard.RUnlock()

	for _, svc := range localPeers {
		if !callback(svc) {
			break
		}
	}
}

func StopLocalPeers() {
	localPeersGuard.RLock()
	defer localPeersGuard.RUnlock()

	for i := len(localPeers) - 1; i >= 0; i-- {
		svc := localPeers[i]
		svc.Stop()
	}
}
