package endpoint

import "reflect"

type Event struct {
	Request interface{}

	Response interface{}
}

type ServiceInfo struct {
	Handler     func(*Event)
	RequestType reflect.Type

	NewResponse func() interface{}
}

type EndPoint interface {
	AddHandler(name string, svc *ServiceInfo)

	// 服务发现注册
	Run() error
}
