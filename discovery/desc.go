package discovery

import "fmt"

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

func (self *ServiceDesc) Address() string {
	return fmt.Sprintf("%s:%d", self.Host, self.Port)
}

func (self *ServiceDesc) String() string {
	return fmt.Sprintf("name: '%s' id: '%s' addr: '%s:%d'  tags: %v", self.Name, self.ID, self.Host, self.Port, self.Tags)
}
