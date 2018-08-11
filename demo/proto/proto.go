package proto

import (
	_ "github.com/davyxu/cellnet/codec/json"
)

func (self *ServerInfo) GetIP() string {
	return self.IP
}

func (self *ServerInfo) GetPort() int {
	return int(self.Port)
}
