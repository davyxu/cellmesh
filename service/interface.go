package service

import (
	"github.com/davyxu/cellnet"
)

type Event interface {
	// 事件对应的会话
	Session() cellnet.Session

	// 事件携带的消息
	Message() interface{}

	GetContextID() int64

	Reply(msg interface{})
}

type DispatcherFunc func(Event)

type Service interface {
	SetDispatcher(dis DispatcherFunc)

	// 服务发现注册
	Start()

	Stop()
}
