package consulsd

import (
	"github.com/hashicorp/consul/api"
	"time"
)

type Config struct {
	*api.Config
	KVCacheDuration      time.Duration // KV缓冲时间，超过时间重新取配置
	ServiceTTL           time.Duration
	ServiceTTLTimeOut    time.Duration
	RemoveServiceTimeout time.Duration
}

func DefaultConfig() *Config {

	return &Config{
		Config:               api.DefaultConfig(),
		KVCacheDuration:      time.Second * 30,            // KV更新间隔
		ServiceTTL:           time.Second * 10,            // 服务心跳
		ServiceTTLTimeOut:    time.Second * 13,            // 服务心跳超时
		RemoveServiceTimeout: time.Minute + 5*time.Second, // Consul移除服务时间，Consul要求必须在1分钟后才能删除TTL超时的服务
	}
}
