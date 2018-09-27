package kvconfig

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/consul"
	"reflect"
)

func doRaw(d discovery.Discovery, key string, defaultValue, ret interface{}) {
	if d == nil {
		return
	}

	err := d.GetValue(key, ret)

	if err == consulsd.ErrValueNotExists {

		reflect.Indirect(reflect.ValueOf(ret)).Set(reflect.ValueOf(defaultValue))
		// 默认值初始化
		d.SetValue(key, defaultValue)
	}

	return
}

// 根据key从Consul中取配置,若不存在,使用默认值且自动写入KV
func String(d discovery.Discovery, key string, defaultValue string) (ret string) {
	doRaw(d, key, defaultValue, &ret)
	return
}

func Int32(d discovery.Discovery, key string, defaultValue int32) (ret int32) {
	doRaw(d, key, defaultValue, &ret)
	return
}

func Int64(d discovery.Discovery, key string, defaultValue int64) (ret int64) {
	doRaw(d, key, defaultValue, &ret)
	return
}

func Bool(d discovery.Discovery, key string, defaultValue bool) (ret bool) {
	doRaw(d, key, defaultValue, &ret)
	return
}
