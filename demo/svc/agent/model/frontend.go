package model

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
)

var (
	FrontendSessionManager peer.SessionManager
	AgentSvcID             string // 网关id
)

func GetClientSession(sesid int64) cellnet.Session {

	return FrontendSessionManager.GetSession(sesid)
}

func GetUser(sesid int64) *User {
	return SessionToUser(GetClientSession(sesid))
}

// 创建一个网关用户
func CreateUser(clientSes cellnet.Session) *User {

	u := NewUser(clientSes)

	// 绑定到session上
	clientSes.(cellnet.ContextSet).SetContext("user", u)
	return u
}

// 用session获取用户
func SessionToUser(clientSes cellnet.Session) *User {

	if clientSes == nil {
		return nil
	}

	if raw, ok := clientSes.(cellnet.ContextSet).GetContext("user"); ok {
		return raw.(*User)
	}

	return nil
}

// 遍历所有的用户
func VisitUser(callback func(*User) bool) {
	FrontendSessionManager.VisitSession(func(clientSes cellnet.Session) bool {

		if u := SessionToUser(clientSes); u != nil {
			return callback(u)
		}

		return true
	})
}
