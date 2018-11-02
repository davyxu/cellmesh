package fxmodel

import (
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"strings"
	"time"
)

type readyChecker interface {
	IsReady() bool
}

func getPeerStatus(svc cellnet.Peer) string {

	type myPeer interface {
		readyChecker
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
		if cs.FetchContext("sd", &desc) {
			context = fmt.Sprintf("--> %22s %22s", desc.ID, desc.Address())
			peerName = desc.Name
		} else {
			context = mp.Address()
			peerName = mp.Name()
		}
	}

	return fmt.Sprintf("%13s %15s %s  [%s]", peerName, mp.TypeName(), context, ready)
}

func LocalServiceStatus() string {

	var sb strings.Builder

	VisitLocalService(func(svc cellnet.Peer) bool {

		if pg, ok := svc.(MultiStatus); ok {

			// 没有连接发现时
			if len(pg.GetPeers()) == 0 {
				sb.WriteString(pg.String())
				sb.WriteString("\n")
			} else {
				for _, p := range pg.GetPeers() {
					sb.WriteString(getPeerStatus(p))
					sb.WriteString("\n")
				}
			}

		} else {
			sb.WriteString(getPeerStatus(svc))
			sb.WriteString("\n")
		}

		return true
	})

	return sb.String()
}

func IsAllReady() (ret bool) {
	ret = true
	VisitLocalService(func(svc cellnet.Peer) bool {
		if !svc.(readyChecker).IsReady() {
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

		time.Sleep(time.Second)

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
