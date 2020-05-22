package redsd

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"reflect"
)

func (self *RedisDiscovery) SetValue(key string, valuePtr interface{}) (ret error) {

	self.operate(func(conn redis.Conn) {
		ret = exec(conn, "SET", key_KV+key, valuePtr)
	})

	return
}

func (self *RedisDiscovery) GetValue(key string, valuePtr interface{}) (ret error) {

	value := reflect.ValueOf(valuePtr)
	if value.Kind() != reflect.Ptr {
		panic("get value must use ptr")
	}

	value = value.Elem()

	self.operate(func(conn redis.Conn) {
		raw, err := conn.Do("GET", key_KV+key)
		ret = err

		if ret != nil {
			return
		}

		switch v := raw.(type) {
		case []byte:

			if value.Kind() == reflect.Struct {
				ret = json.Unmarshal(v, valuePtr)
			} else {
				value.SetBytes(v)
			}

		case string:
			value.SetString(v)
		case int64:
			value.SetInt(v)
		default:
			value.Set(reflect.ValueOf(v))
		}
	})

	return
}
