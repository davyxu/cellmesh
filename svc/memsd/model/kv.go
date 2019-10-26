package model

import (
	"encoding/json"
	"github.com/davyxu/cellmesh/discovery"
	"io"
	"sort"
)

type ValueMeta struct {
	Key     string
	Value   []byte
	SvcName string // 服务才有此名字
	Token   string
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

	ValueDirty bool
)

func SetValue(key string, meta *ValueMeta) {
	ValueDirty = true
	valueByKey[key] = meta
}

func GetValue(key string) *ValueMeta {

	return valueByKey[key]
}

func DeleteValue(key string) *ValueMeta {
	ValueDirty = true
	ret := valueByKey[key]
	delete(valueByKey, key)

	return ret
}

func ValueCount() int {
	return len(valueByKey)
}

func VisitValue(callback func(*ValueMeta) bool) {
	for _, vmeta := range valueByKey {
		if !callback(vmeta) {
			return
		}
	}
}

type PersistFile struct {
	Version int
	Values  []*ValueMeta
}

var (
	fileVersion = 1
)

func SaveValue(writer io.Writer) (valuesSaved int, err error) {

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")

	var file PersistFile
	file.Version = fileVersion
	for _, vmeta := range valueByKey {

		// 服务不保存, 服务重新注册
		if IsServiceKey(vmeta.Key) {
			continue
		}

		file.Values = append(file.Values, vmeta)
	}

	sort.SliceStable(file.Values, func(i, j int) bool {

		return file.Values[i].Key < file.Values[j].Key
	})

	err = encoder.Encode(&file)

	if err != nil {
		return
	}

	return len(file.Values), nil
}

func LoadValue(reader io.Reader) error {

	decoder := json.NewDecoder(reader)

	var file PersistFile
	err := decoder.Decode(&file)
	if err != nil {
		return err
	}

	valueByKey = map[string]*ValueMeta{}

	for _, v := range file.Values {
		valueByKey[v.Key] = v
	}

	return nil
}
