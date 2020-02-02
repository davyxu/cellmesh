package fx

import (
	"reflect"
	"sync"
)

type InjectInvoker func(ioc *InjectContext) interface{}

type InjectContext struct {
	invokerByKey      map[interface{}]InjectInvoker
	invokerByKeyGuard sync.RWMutex
	parent            *InjectContext
}

func (self *InjectContext) SetParent(p *InjectContext) {
	self.parent = p
}

func (self *InjectContext) MapFunc(key interface{}, invoker InjectInvoker) {

	if self.findType(key) != nil {
		panic("duplicate MapFunc key")
	}

	self.invokerByKeyGuard.Lock()
	self.invokerByKey[key] = invoker
	self.invokerByKeyGuard.Unlock()
}

func (self *InjectContext) findType(key interface{}) InjectInvoker {

	self.invokerByKeyGuard.RLock()
	if v, ok := self.invokerByKey[key]; ok {
		self.invokerByKeyGuard.RUnlock()
		return v
	}
	self.invokerByKeyGuard.RUnlock()

	if self.parent != nil {
		return self.parent.findType(key)
	}

	return nil
}

func (self *InjectContext) TryInvoke(key interface{}) (value interface{}, ok bool) {
	invoker := self.findType(key)
	if invoker == nil {
		return nil, false
	}

	return invoker(self), true
}

func keyToString(key interface{}) string {
	switch v := key.(type) {
	case string:
		return v
	case reflect.Type:
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		return v.Name()
	default:

		t := reflect.TypeOf(key)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		return t.Name()
	}

}

func (self *InjectContext) Invoke(key interface{}) interface{} {

	if value, ok := self.TryInvoke(key); ok {
		return value
	} else {
		panic("type not register mapper: " + keyToString(key))
	}
}

func NewInjectContext() *InjectContext {
	return &InjectContext{
		invokerByKey: make(map[interface{}]InjectInvoker),
	}

}
