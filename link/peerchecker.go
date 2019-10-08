package link

import (
	"fmt"
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"strings"
	"time"
)

func peerStatus(svc cellnet.Peer) string {

	type myPeer interface {
		cellnet.PeerReadyChecker
		Name() string
		Address() string
		cellnet.Peer
	}

	mp := svc.(myPeer)

	var ready string
	if mp.IsReady() {
		ready = "READY"
	}

	var peerName string
	var context string
	if cs, ok := svc.(cellnet.ContextSet); ok {

		var desc *discovery.ServiceDesc
		if cs.FetchContext(cellmesh.PeerContextKey_ServiceDesc, &desc) {
			context = fmt.Sprintf("  %22s %22s", desc.ID, desc.Address())
			peerName = desc.Name
		} else {
			if sesGetter, ok := svc.(interface {
				Session() cellnet.Session
			}); ok {

				ses := sesGetter.Session()

				svcID := GetRemoteLinkSvcID(ses)

				context = fmt.Sprintf("  %22s %22s", svcID, mp.Address())
			} else {
				context = mp.Address()

			}

			peerName = mp.Name()

		}
	}

	return fmt.Sprintf("%13s %15s %s  [%s]", peerName, mp.TypeName(), context, ready)
}

func remoteSessionStatus(ses cellnet.Session) string {
	svcID := GetRemoteLinkSvcID(ses)
	svcName := GetRemoteLinkSvcName(ses)

	return fmt.Sprintf("%s(%s)", svcName, svcID)
}

func LocalServiceStatus() string {

	var sb strings.Builder

	VisitLocalPeer(func(svc cellnet.Peer) bool {

		sb.WriteString(peerStatus(svc))
		sb.WriteString("\n")

		return true
	})

	//VisitRemoteSession(func(ses cellnet.Session) bool {
	//	sb.WriteString(remoteSessionStatus(ses))
	///	sb.WriteString("\n")
	//	return true
	//})

	return sb.String()
}

func IsAllReady() (ret bool) {
	ret = true
	VisitLocalPeer(func(svc cellnet.Peer) bool {
		if !svc.(cellnet.PeerReadyChecker).IsReady() {
			ret = false
			return false
		}

		return true
	})

	return
}

func CheckReady() {

	var lastStatus string
	for {

		time.Sleep(time.Second * 3)

		if IsAllReady() {
			log.SetColor("green").Infof("All peers ready!\n%s", LocalServiceStatus())

			break
		}

		thisStatus := LocalServiceStatus()

		if lastStatus != thisStatus {
			log.Warnf("peers not all ready\n%s", thisStatus)
			lastStatus = thisStatus
		}

	}

}
