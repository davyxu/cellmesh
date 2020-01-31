package fx

import (
	"fmt"
	"testing"
)

type User struct {
	ID int
}

type Message struct {
	MsgID  int
	UserID int
}

func GetUser(id int) *User {
	return &User{
		ID: id,
	}
}

func GetCallback(msgid int) interface{} {

	switch msgid {
	case 1:
		return MessageWithUserFunc(func(ioc *InjectContext, msg *Message, u *User) {
			fmt.Println(msg, u)
		})
	case 2:
		return MessageWithUserIDFunc(func(ioc *InjectContext, msg *Message, userID int) {
			fmt.Println(msg, userID)
		})
	}

	return nil
}

var globalIOC *InjectContext

func init() {
	globalIOC = NewInjectContext()

	globalIOC.MapFunc("User", func(ioc *InjectContext) interface{} {
		msg := ioc.Invoke("Message").(*Message)

		return GetUser(msg.UserID)
	})

	globalIOC.MapFunc("CallHandler", func(ioc *InjectContext) interface{} {
		msg := ioc.Invoke("Message").(*Message)

		u := ioc.Invoke("User").(*User)

		callback := GetCallback(msg.MsgID)
		if callback == nil {
			return nil
		}

		switch f := callback.(type) {
		case MessageWithUserFunc:
			f(ioc, msg, u)
		case MessageWithUserIDFunc:
			f(ioc, msg, msg.UserID)
		}

		return nil
	})
}

func genIOCConext(msg *Message) *InjectContext {
	// 框架层
	ioc := NewInjectContext()

	ioc.SetParent(globalIOC)

	ioc.MapFunc("Message", func(ioc *InjectContext) interface{} {
		return msg
	})

	return ioc
}

func TestIOC(t *testing.T) {

	ioc := genIOCConext(&Message{UserID: 10, MsgID: 1})
	OnMessage(ioc)

	ioc = genIOCConext(&Message{UserID: 10, MsgID: 2})
	OnMessage(ioc)
}

type MessageWithUserFunc func(ioc *InjectContext, msg *Message, u *User)

type MessageWithUserIDFunc func(ioc *InjectContext, msg *Message, useID int)

func OnMessage(ioc *InjectContext) {

	ioc.Invoke("CallHandler")

}
