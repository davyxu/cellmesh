package fx

import (
	"github.com/davyxu/cellnet"
	"reflect"
)

var (
	MessageRegistry = NewInjectContext()
)

func RegisterMessage(msgTypeObj interface{}, handler func(ioc *InjectContext, ev cellnet.Event)) {

	tMsg := reflect.TypeOf(msgTypeObj)
	if tMsg.Kind() != reflect.Ptr {
		panic("require msg ptr")
	}

	MessageRegistry.MapFunc(tMsg.Elem(), func(ioc *InjectContext) interface{} {
		ev := ioc.Invoke("Event").(cellnet.Event)
		handler(ioc, ev)

		return nil
	})

}
