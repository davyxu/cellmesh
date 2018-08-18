package backend

import (
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
)

// 将客户端连接绑定到后台服务
func bindClientToBackend(backendSes cellnet.Session, clientSesID int64) {

	clientSes := model.GetClientSession(clientSesID)
	if clientSes == nil {
		return
	}

	sd := service.ServiceSessionToSD(backendSes)
	if sd == nil {
		log.Errorln("backend sd not found")
		return
	}

	u := model.SessionToUser(clientSes)

	if u != nil {
		u.SetBackend(sd.Name, backendSes)
	} else {
		u = model.CreateUser(clientSes)
		u.AddBackend(sd.Name, backendSes)
	}

}

// 恢复后台连接
func recoverBackend(backendSes cellnet.Session, svcName string) {

	model.VisitUser(func(u *model.User) bool {
		u.SetBackend(svcName, backendSes)

		return true
	})

}

// 移除玩家对应的后台连接
func removeBackend(backendSes cellnet.Session) {

	sd := service.ServiceSessionToSD(backendSes)
	if sd == nil {
		log.Errorln("backend sd not found")
		return
	}

	model.VisitUser(func(u *model.User) bool {
		u.SetBackend(sd.Name, nil)

		return true
	})

}
