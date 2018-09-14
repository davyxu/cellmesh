package discovery

import (
	"fmt"
)

type ServiceDesc struct {
	Name string
	ID   string // 所有service中唯一的id
	Host string
	Port int
	Tags []string // 标签
}

func (self *ServiceDesc) Address() string {
	return fmt.Sprintf("%s:%d", self.Host, self.Port)
}

func (self *ServiceDesc) String() string {
	return fmt.Sprintf("name: '%s' id: '%s' addr: '%s:%d'", self.Name, self.ID, self.Host, self.Port)
}

type Discovery interface {

	// 注册服务
	Register(*ServiceDesc) error

	// 解注册服务
	Deregister(svcid string) error

	// 根据服务名查到可用的服务
	Query(name string) (ret []*ServiceDesc, err error)

	// 注册服务变化通知
	RegisterNotify(mode string) (ret chan struct{})

	// 解除服务变化通知
	DeregisterNotify(mode string, c chan struct{})

	// https://www.consul.io/intro/getting-started/kv.html
	// 设置值
	SetValue(key string, value []byte) error

	// 获取值
	GetValue(key string) ([]byte, bool, error)
}

var (
	Default Discovery
)
