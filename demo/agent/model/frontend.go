package model

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
)

var (
	FrontendListener cellnet.Peer
)

func GetClientSession(sesid int64) cellnet.Session {

	return FrontendListener.(peer.SessionManager).GetSession(sesid)
}

// 创建一个网关用户
func CreateUser(clientSes cellnet.Session) *User {

	u := NewUser()

	// 绑定到session上
	clientSes.(cellnet.ContextSet).SetContext("user", u)
	return u
}

// 用session获取用户
func GetUser(clientSes cellnet.Session) *User {

	if raw, ok := clientSes.(cellnet.ContextSet).GetContext("user"); ok {
		return raw.(*User)
	}

	return nil
}

// 遍历所有的用户
func VisitUser(callback func(*User) bool) {
	FrontendListener.(peer.SessionManager).VisitSession(func(clientSes cellnet.Session) bool {

		if u := GetUser(clientSes); u != nil {
			return callback(u)
		}

		return true
	})
}
