package fx

import (
	"github.com/davyxu/cellnet"
	"reflect"
)

func MakeIOCEventHandler(parentIOC *InjectContext) cellnet.EventCallback {

	return func(ev cellnet.Event) {
		// 框架层
		ioc := NewInjectContext()

		ioc.SetParent(parentIOC)

		ioc.MapFunc("Event", func(ioc *InjectContext) interface{} {
			return ev
		})

		tMsg := reflect.TypeOf(ev.Message())
		if tMsg.Kind() == reflect.Ptr {
			tMsg = tMsg.Elem()
		}
		ioc.TryInvoke(tMsg)
	}
}
