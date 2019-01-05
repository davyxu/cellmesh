package discovery

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// 注册到服务发现的服务描述
type ServiceDesc struct {
	Name string
	ID   string // 所有service中唯一的id
	Host string
	Port int
	Tags []string          // 分类标签
	Meta map[string]string // 细节配置
}

func (self *ServiceDesc) Equals(sd *ServiceDesc) bool {

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

	if !reflect.DeepEqual(self.Tags, sd.Tags) {
		return false
	}

	if !reflect.DeepEqual(self.Meta, sd.Meta) {
		return false
	}

	return true
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

	return fmt.Sprintf("%s host: %s port: %d %s", self.ID, self.Host, self.Port, sb.String())
}

func (self *ServiceDesc) FormatString() string {

	var sb strings.Builder
	if len(self.Meta) > 0 {

		type pair struct {
			key   string
			value string
		}

		var pairs []pair

		for key, value := range self.Meta {
			pairs = append(pairs, pair{key, value})
		}

		sort.Slice(pairs, func(i, j int) bool {

			return pairs[i].key < pairs[j].key
		})

		sb.WriteString("meta: [ ")
		for _, kv := range pairs {
			sb.WriteString(kv.key)
			sb.WriteString("=")
			sb.WriteString(kv.value)
			sb.WriteString(" ")
		}
		sb.WriteString("]")
	}

	return fmt.Sprintf("%25s host: %15s port: %5d %s", self.ID, self.Host, self.Port, sb.String())
}
