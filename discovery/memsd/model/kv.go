package model

import (
	"encoding/json"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
)

type ValueMeta struct {
	Key     string
	Value   []byte
	SvcName string          // 服务才有此名字
	Ses     cellnet.Session // 来源的连接
}

var ErrDesc = discovery.ServiceDesc{Name: "invalid desc"}

func (self *ValueMeta) ValueAsServiceDesc() *discovery.ServiceDesc {

	var desc discovery.ServiceDesc
	err := json.Unmarshal(self.Value, &desc)
	if err != nil {
		return &ErrDesc
	}

	return &desc
}

var (
	valueByKey = map[string]*ValueMeta{}
)

func SetValue(key string, meta *ValueMeta) {

	valueByKey[key] = meta
}

func GetValue(key string) *ValueMeta {

	return valueByKey[key]
}

func DeleteValue(key string) *ValueMeta {
	ret := valueByKey[key]
	delete(valueByKey, key)

	return ret
}

func VisitValue(callback func(*ValueMeta) bool) {
	for _, vmeta := range valueByKey {
		if !callback(vmeta) {
			return
		}
	}
}
