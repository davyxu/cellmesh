package discovery

import (
	"fmt"
)

type ServiceDesc struct {
	Name    string
	ID      string
	Address string
	Port    int
}

func (self *ServiceDesc) String() string {
	return fmt.Sprintf("name: '%s' id: '%s' addr: '%s:%d'", self.Name, self.ID, self.Address, self.Port)
}

type NotifyFunc func(eventName string, desc *ServiceDesc)

type Discovery interface {

	// 注册服务
	Register(*ServiceDesc) error

	// 解注册服务
	Deregister(svcid string) error

	// 根据服务名查到可用的服务
	Query(name string) (ret []*ServiceDesc, err error)

	// 提供一个版本号，只有当版本变化时才停止阻塞
	WaitAdded()
}

var (
	Default Discovery
)
