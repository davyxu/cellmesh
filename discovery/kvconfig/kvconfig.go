package kvconfig

import (
	"github.com/davyxu/cellmesh/discovery"
	"strconv"
)

// 根据key从Consul中取配置,若不存在,使用默认值且自动写入KV
func String(key string, defaultValue string) (ret string) {

	data, exists, err := discovery.Default.GetValue(key)
	if err != nil {
		ret = defaultValue
		return
	}

	v := string(data)

	if !exists {
		ret = defaultValue

		// 默认值初始化
		defer discovery.Default.SetValue(key, []byte(ret))
		return
	}

	ret = v

	return
}

func Int32(key string, defaultValue int32) (ret int32) {

	data, exists, err := discovery.Default.GetValue(key)
	if err != nil {
		ret = defaultValue
		return
	}

	v, err := strconv.ParseInt(string(data), 10, 32)
	if err != nil {
		ret = defaultValue
	}

	if !exists {
		ret = defaultValue

		// 默认值初始化
		defer discovery.Default.SetValue(key, []byte(strconv.FormatInt(int64(ret), 10)))
		return
	}

	ret = int32(v)

	return
}

func Int64(key string, defaultValue int64) (ret int64) {

	data, exists, err := discovery.Default.GetValue(key)
	if err != nil {
		ret = defaultValue
		return
	}

	v, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		ret = defaultValue
	}

	if !exists {
		ret = defaultValue

		// 默认值初始化
		defer discovery.Default.SetValue(key, []byte(strconv.FormatInt(ret, 10)))
		return
	}

	ret = v

	return
}

func Bool(key string, defaultValue bool) (ret bool) {

	data, exists, err := discovery.Default.GetValue(key)
	if err != nil {
		ret = defaultValue
		return
	}

	v, err := strconv.ParseBool(string(data))
	if err != nil {
		ret = defaultValue
	}

	if !exists {
		ret = defaultValue

		// 默认值初始化
		defer discovery.Default.SetValue(key, []byte(strconv.FormatBool(ret)))
		return
	}

	ret = v

	return
}
