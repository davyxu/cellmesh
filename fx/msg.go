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

		msgType := reflect.TypeOf(ev.Message())

		ioc.Invoke(msgType)
	}
}
