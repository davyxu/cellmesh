package svc

type Event struct {
	Request interface{}

	Response interface{}
}

type Handler func(*Event)

type Service interface {
	AddHandler(name string, handler Handler)

	// 服务发现注册
	Run() error
}
