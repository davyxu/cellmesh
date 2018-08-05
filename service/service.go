package service

import "reflect"

type Event struct {
	Request interface{}

	Response interface{}
}

type MethodInfo struct {
	Handler     func(*Event)
	RequestType reflect.Type

	NewResponse func() interface{}
}

type Service interface {
	AddMethod(name string, svc *MethodInfo)

	// 服务发现注册
	Run() error
}

var (
	NewService func(name string) Service
)
