package fxmodel

import (
	"github.com/davyxu/cellnet"
	"sync"
)

var (
	localServices      []cellnet.Peer
	localServicesGuard sync.RWMutex
)

func AddLocalService(p cellnet.Peer) {
	localServicesGuard.Lock()
	localServices = append(localServices, p)
	localServicesGuard.Unlock()
}

func RemoveLocalService(p cellnet.Peer) {
	localServicesGuard.Lock()
	for index, libp := range localServices {
		if libp == p {
			localServices = append(localServices[:index], localServices[index+1:]...)
			break
		}

	}
	localServicesGuard.Unlock()
}

func GetLocalService(svcName string) cellnet.Peer {

	localServicesGuard.RLock()
	defer localServicesGuard.RUnlock()

	for _, svc := range localServices {
		if prop, ok := svc.(cellnet.PeerProperty); ok && prop.Name() == svcName {
			return svc
		}
	}

	return nil
}

func VisitLocalService(callback func(cellnet.Peer) bool) {
	localServicesGuard.RLock()
	defer localServicesGuard.RUnlock()

	for _, svc := range localServices {
		if !callback(svc) {
			break
		}
	}
}

func StopAllService() {
	localServicesGuard.RLock()
	defer localServicesGuard.RUnlock()

	for i := len(localServices) - 1; i >= 0; i-- {
		svc := localServices[i]
		svc.Stop()
	}
}
