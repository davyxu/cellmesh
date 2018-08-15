package service

import (
	"github.com/davyxu/cellnet"
	"reflect"
)

type Event struct {
	Session cellnet.Session
	Request interface{}

	Response interface{}
}

type MethodInfo struct {
	Handler     func(*Event)
	RequestType reflect.Type

	NewResponse func() interface{}
}

type Service interface {
	AddCall(name string, svc *MethodInfo)

	// 服务发现注册
	Run() error

	ID() string
}
