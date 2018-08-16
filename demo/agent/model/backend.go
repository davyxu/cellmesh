package model

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
)

func CreateUser(clientSes cellnet.Session) *User {

	u := NewUser()

	clientSes.(cellnet.ContextSet).SetContext("user", u)
	return u
}

func GetUser(clientSes cellnet.Session) *User {

	if raw, ok := clientSes.(cellnet.ContextSet).GetContext("user"); ok {
		return raw.(*User)
	}

	return nil
}

func VisitUser(callback func(*User) bool) {
	FrontendListener.(peer.SessionManager).VisitSession(func(clientSes cellnet.Session) bool {

		if u := GetUser(clientSes); u != nil {
			return callback(u)
		}

		return true
	})
}
