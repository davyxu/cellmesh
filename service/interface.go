package service

import (
	"github.com/davyxu/cellnet"
)

type Event interface {
	// 给来源会话(网关,服务)发消息
	Session() cellnet.Session

	// 事件携带的消息
	Message() interface{}

	// 网关透传输出,如客户端在网关的SessionID
	PassThrough() interface{}

	// 回复客户端
	Reply(msg interface{})
}

type DispatcherFunc func(Event)

type Service interface {
	SetDispatcher(dis DispatcherFunc)

	// 服务发现注册
	Start()

	Stop()
}
