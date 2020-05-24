package backend

import (
	"errors"
	"github.com/davyxu/cellmesh/svc/node/agent/model"
	"github.com/davyxu/ulog"
)

var (
	ErrAlreadyBind           = errors.New("already bind user")
	ErrBackendServerNotFound = errors.New("backend svc not found")
)

// 将客户端连接绑定到后台服务
func bindClientToBackend(nodeID string, clientSesID int64) {
	// 将客户端的id转为session
	clientSes := model.GetClientSession(clientSesID)

	if clientSes == nil {
		return
	}

	// 从客户端的会话取得用户
	u := model.SessionToUser(clientSes)

	// 已经绑定
	if u != nil {
		ulog.Warnf("duplicate user bind backend, nodeid: %s, sesID: %d", nodeID, clientSesID)
		return
	}

	u = model.CreateUser(clientSes)

	// 更新绑定后台服务的svcid
	u.BindBackend(nodeID)
}
