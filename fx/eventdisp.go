package fx

import (
	"github.com/davyxu/cellnet"
	"reflect"
)

var (
	MessageRegistry = NewInjectContext()
)

func RegisterMessage(msgTypeObj interface{}, handler func(ioc *InjectContext, ev cellnet.Event)) {
	msgType := reflect.TypeOf(msgTypeObj)

	MessageRegistry.MapFunc(msgType, func(ioc *InjectContext) interface{} {
		ev := ioc.Invoke("Event").(cellnet.Event)
		handler(ioc, ev)

		return nil
	})

}
