package consulsd

import (
	"github.com/hashicorp/consul/api"
	"time"
)

type Config struct {
	*api.Config
	ServiceTTL           time.Duration
	ServiceTTLTimeOut    time.Duration
	RemoveServiceTimeout time.Duration
}

func DefaultConfig() *Config {

	return &Config{
		Config:               api.DefaultConfig(),
		ServiceTTL:           time.Second * 10,            // 服务心跳
		ServiceTTLTimeOut:    time.Second * 20,            // 服务心跳超时(注意，超时时间太短时，会导致服务器不会出现在列表中)
		RemoveServiceTimeout: time.Minute + 5*time.Second, // Consul移除服务时间，Consul要求必须在1分钟后才能删除TTL超时的服务
	}
}
