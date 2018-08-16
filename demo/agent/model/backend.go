package model

import (
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
)

func BindClientToBackend(backendSes cellnet.Session, clientSesID int64) {

	clientSes := GetClient(clientSesID)
	if clientSes == nil {
		return
	}

	sd := service.GetSessionSD(backendSes)
	if sd == nil {
		log.Errorln("backend sd not found")
		return
	}

	clientSes.(cellnet.ContextSet).SetContext("route_"+sd.Name, backendSes)
}

func GetClientBackendSession(clientSes cellnet.Session, svcName string) cellnet.Session {

	if raw, ok := clientSes.(cellnet.ContextSet).GetContext("route_" + svcName); ok {
		return raw.(cellnet.Session)
	}

	return nil
}
