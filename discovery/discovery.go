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
	return fmt.Sprintf("name: %s id: %s addr: %s:%d", self.Name, self.ID, self.Address, self.Port)
}

type Discovery interface {
	Register(*ServiceDesc) error

	Deregister(svcid string) error

	Query(name string) (ret []*ServiceDesc, err error)
}

var (
	Default Discovery
)
