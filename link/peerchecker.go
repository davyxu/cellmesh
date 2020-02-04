package link

import (
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/ulog"
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
		ready = "[READY]"
	}

	var peerName string
	var context string
	if cs, ok := svc.(cellnet.ContextSet); ok {

		var desc *discovery.ServiceDesc
		if cs.FetchContext(PeerContextKey_ServiceDesc, &desc) {
			context = fmt.Sprintf("  %22s %22s", desc.ID, desc.Address())
			peerName = desc.Name
		} else {
			if sesGetter, ok := svc.(interface {
				Session() cellnet.Session
			}); ok {

				ses := sesGetter.Session()

				svcID := GetLinkSvcID(ses)

				context = fmt.Sprintf("  %22s %22s", svcID, mp.Address())
			} else {
				context = mp.Address()

			}

			peerName = mp.Name()

		}
	} else {
		peerName = mp.Name()
	}

	return fmt.Sprintf("%13s %15s %s  %s", peerName, mp.TypeName(), context, ready)
}

func LocalServiceStatus() string {

	var sb strings.Builder

	VisitPeer(func(svc cellnet.Peer) bool {

		sb.WriteString(peerStatus(svc))
		sb.WriteString("\n")

		return true
	})

	return sb.String()
}

func IsAllReady() (ret bool) {
	ret = true
	VisitPeer(func(svc cellnet.Peer) bool {
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
			ulog.WithColorName("green").Infof("All peers ready!\n%s", LocalServiceStatus())

			break
		}

		thisStatus := LocalServiceStatus()

		if lastStatus != thisStatus {
			ulog.Warnf("peers not all ready\n%s", thisStatus)
			lastStatus = thisStatus
		}

	}

}
