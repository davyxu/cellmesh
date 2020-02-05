package discovery

type ValueMeta struct {
	Key   string
	Value []byte
}

type NotifyContext struct {
	Mode string
	Desc *ServiceDesc
}

// 基础服务发现
type Discovery interface {
	Start(config interface{})

	// 注册服务
	Register(*ServiceDesc) error

	// 解注册服务
	Deregister(svcid string) error

	// 根据服务名查到可用的服务
	Query(name string) (ret []*ServiceDesc)

	// 注册服务变化通知
	// 'add'表示有服务加入
	// 'ready'表示服务发现准备好
	// 'mod' 表示有服务修改
	// 'del' 表示有服务删除
	RegisterNotify() (ret chan *NotifyContext)

	// 解除服务变化通知
	DeregisterNotify(c chan struct{})

	// 设置值
	SetValue(key string, value interface{}, optList ...interface{}) error

	// 取值，并赋值到变量
	GetValue(key string, valuePtr interface{}) error

	// 删除值
	DeleteValue(key string) error
}

var (
	Global Discovery
)
