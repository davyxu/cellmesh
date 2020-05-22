package frontend

import (
	"errors"
	"github.com/davyxu/cellmesh/link"
	"github.com/davyxu/cellmesh/svc/agent/model"
)

var (
	ErrAlreadyBind           = errors.New("already bind user")
	ErrBackendServerNotFound = errors.New("backend svc not found")
)

// 将客户端连接绑定到后台服务
func bindClientToBackend(backendSvcID string, clientSesID int64) (*model.User, error) {

	desc := link.DescByID(backendSvcID)

	if desc == nil {
		return nil, ErrBackendServerNotFound
	}

	// 将客户端的id转为session
	clientSes := model.GetClientSession(clientSesID)

	// 从客户端的会话取得用户
	u := model.SessionToUser(clientSes)

	// 已经绑定
	if u != nil {
		return nil, ErrAlreadyBind
	}

	u = model.CreateUser(clientSes)

	// 更新绑定后台服务的svcid
	u.BindBackend(desc.Name, desc.ID)

	return u, nil
}
