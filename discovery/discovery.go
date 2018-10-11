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
	SetValue(key string, value interface{}) error

	GetValue(key string, valuePtr interface{}) error

	// 获取原始值
	GetRawValue(key string) ([]byte, error)

	// 获取原始值列表
	GetRawValueList(key string) ([]ValueMeta, error)

	// 设置服务状态汇报
	SetChecker(svcid string, checkerFunc CheckerFunc)
}

var (
	Default Discovery
)
