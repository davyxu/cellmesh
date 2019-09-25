package discovery

type ValueMeta struct {
	Key   string
	Value []byte
}

type NotifyFunc func(evType string, args ...interface{})

// 基础服务发现
type Discovery interface {
	Start(config interface{})

	// 注册服务
	Register(*ServiceDesc) error

	// 解注册服务
	Deregister(svcid string) error

	// 根据服务名查到可用的服务
	Query(name string) (ret []*ServiceDesc)

	// 在服务发现内部逻辑线程, 通知回调, 多个事件将顺序通知
	SetNotify(callback NotifyFunc)
}

// KV接口, 可由Discovery转换
type KVStorage interface {
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
