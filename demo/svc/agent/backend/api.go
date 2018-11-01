package backend

import (
	"github.com/davyxu/cellmesh/demo/svc/agent/model"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
)

// 将客户端连接绑定到后台服务
func bindClientToBackend(backendSes cellnet.Session, clientSesID int64) {

	// 将客户端的id转为session
	clientSes := model.GetClientSession(clientSesID)

	// 客户端已经断开了
	if clientSes == nil {
		return
	}

	// 取得后台服务的信息
	sd := service.SessionToContext(backendSes)
	if sd == nil {
		log.Errorln("backend sd not found")
		return
	}

	// 从客户端的会话取得用户
	u := model.SessionToUser(clientSes)

	// 第一次绑定
	if u == nil {
		u = model.CreateUser(clientSes)
	}

	// 更新绑定后台服务的svcid
	u.SetBackend(sd.Name, sd.SvcID)

}
