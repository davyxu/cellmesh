package service

import (
	"github.com/davyxu/cellnet"
	"reflect"
)

type Event struct {
	Session cellnet.Session

	Request interface{}

	Response interface{}

	ContextID []int64
}

type MethodInfo struct {
	Handler     func(*Event)
	RequestType reflect.Type

	NewResponse func() interface{}
}

type Service interface {
	SetDispatcher(dis *Dispatcher)

	// 服务发现注册
	Start()

	Stop()
}

type ReplyEvent interface {
	Reply(msg interface{})
}
