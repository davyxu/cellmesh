package redsd

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"reflect"
	"strconv"
	"strings"
)

// 注册到服务发现的服务描述
type NodeDesc struct {
	Name    string            // 进程类型
	ID      string            // 唯一的id
	Host    string            // 地址
	Port    int               // 端口
	Meta    map[string]string // 细节配置
	Peer    cellnet.Peer      `json:"-"` // 连接时无, 侦听时有
	Session cellnet.Session   `json:"-"` // Accept连接上来的有此Session
}

func (self *NodeDesc) Equals(sd *NodeDesc) bool {

	if sd.ID != self.ID {
		return false
	}

	if sd.Port != self.Port {
		return false
	}

	if sd.Name != self.Name {
		return false
	}

	if sd.Host != self.Host {
		return false
	}

	if !reflect.DeepEqual(self.Meta, sd.Meta) {
		return false
	}

	return true
}

func (self *NodeDesc) SetMeta(key, value string) {
	if self.Meta == nil {
		panic("meta not init")
	}

	self.Meta[key] = value
}

func (self *NodeDesc) GetMeta(name string) string {
	if self.Meta == nil {
		panic("meta not init")
	}

	return self.Meta[name]
}

func (self *NodeDesc) GetMetaAsInt(name string) int {
	v, err := strconv.ParseInt(self.GetMeta(name), 10, 64)
	if err != nil {
		return 0
	}

	return int(v)
}

func (self *NodeDesc) Address() string {

	if self.Host == "" && self.Port == 0 {
		return ""
	}

	return fmt.Sprintf("%s:%d", self.Host, self.Port)
}

func (self *NodeDesc) String() string {
	var sb strings.Builder
	if len(self.Meta) > 0 {

		sb.WriteString("meta: [ ")
		for key, value := range self.Meta {
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(value)
			sb.WriteString(" ")
		}
		sb.WriteString("]")
	}

	return fmt.Sprintf("%s host: %s port: %d %s", self.ID, self.Host, self.Port, sb.String())
}

func NewDesc() *NodeDesc {
	return &NodeDesc{
		Meta: make(map[string]string),
	}
}
