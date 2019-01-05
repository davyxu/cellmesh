package discovery

type ValueMeta struct {
	Key   string
	Value []byte
}

type CheckerFunc func() (output, status string)

type Discovery interface {

	// 注册服务
	Register(*ServiceDesc) error

	// 解注册服务
	Deregister(svcid string) error

	// 根据服务名查到可用的服务
	Query(name string) (ret []*ServiceDesc)

	// 注册服务变化通知
	RegisterNotify(mode string) (ret chan struct{})

	// 解除服务变化通知
	DeregisterNotify(mode string, c chan struct{})

	// https://www.consul.io/intro/getting-started/kv.html
	// 设置值
	SetValue(key string, value interface{}, optList ...interface{}) error

	// 取值，并赋值到变量
	GetValue(key string, valuePtr interface{}) error

	// 删除值
	DeleteValue(key string) error
}

var (
	Default Discovery
)
