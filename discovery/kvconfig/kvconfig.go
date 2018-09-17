package kvconfig

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/consul"
)

var (
	di discovery.Discovery
)

// 设置自定义的接口
func SetDiscovery(d discovery.Discovery) {
	di = d
}

func getDI() discovery.Discovery {

	if di == nil {
		return discovery.Default
	}

	return di
}

func doRaw(key string, defaultValue, ret interface{}) {
	err := getDI().GetValue(key, ret)

	if err == consulsd.ErrValueNotExists {

		ret = defaultValue
		// 默认值初始化
		getDI().SetValue(key, defaultValue)
	}

	return
}

// 根据key从Consul中取配置,若不存在,使用默认值且自动写入KV
func String(key string, defaultValue string) (ret string) {
	doRaw(key, defaultValue, &ret)
	return
}

func Int32(key string, defaultValue int32) (ret int32) {
	doRaw(key, defaultValue, &ret)
	return
}

func Int64(key string, defaultValue int64) (ret int64) {
	doRaw(key, defaultValue, &ret)
	return
}

func Bool(key string, defaultValue bool) (ret bool) {
	doRaw(key, defaultValue, &ret)
	return
}
