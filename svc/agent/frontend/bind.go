package frontend

import (
	"errors"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh_demo/svc/agent/model"
)

var (
	ErrAlreadyBind           = errors.New("already bind user")
	ErrBackendSDNotFound     = errors.New("backend sd not found")
	ErrBackendServerNotFound = errors.New("backend svc not found")
)

// 将客户端连接绑定到后台服务
func bindClientToBackend(backendSvcID string, clientSesID int64) (*model.User, error) {

	backendSes := service.GetRemoteService(backendSvcID)

	if backendSes == nil {
		return nil, ErrBackendServerNotFound
	}

	// 取得后台服务的信息
	sd := service.SessionToContext(backendSes)
	if sd == nil {
		return nil, ErrBackendSDNotFound
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
	u.SetBackend(sd.Name, sd.SvcID)

	return u, nil
}
