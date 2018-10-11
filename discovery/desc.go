package discovery

import (
	"fmt"
	"strconv"
	"strings"
)

type ServiceDesc struct {
	Name string
	ID   string // 所有service中唯一的id
	Host string
	Port int
	Tags []string // 标签
	Meta map[string]string
}

func (self *ServiceDesc) ContainTags(tag string) bool {
	for _, libtag := range self.Tags {
		if libtag == tag {
			return true
		}
	}

	return false
}

func (self *ServiceDesc) SetMeta(key, value string) {
	if self.Meta == nil {
		self.Meta = make(map[string]string)
	}

	self.Meta[key] = value
}

func (self *ServiceDesc) GetMeta(name string) string {
	if self.Meta == nil {
		return ""
	}

	return self.Meta[name]
}

func (self *ServiceDesc) GetMetaAsInt(name string) int {
	v, err := strconv.ParseInt(self.GetMeta(name), 10, 64)
	if err != nil {
		return 0
	}

	return int(v)
}

func (self *ServiceDesc) Address() string {
	return fmt.Sprintf("%s:%d", self.Host, self.Port)
}

func (self *ServiceDesc) String() string {

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

	return fmt.Sprintf("%15s port: %5d %s", self.ID, self.Port, sb.String())
}
