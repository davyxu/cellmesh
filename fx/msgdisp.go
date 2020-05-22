package fx

import (
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"reflect"
)

var (
	MessageRegistry = meshutil.NewInjectContext()
)

func RegisterMessage(msgTypeObj interface{}, handler func(ioc *meshutil.InjectContext, ev cellnet.Event)) {

	tMsg := reflect.TypeOf(msgTypeObj)
	if tMsg.Kind() != reflect.Ptr {
		panic("require msg ptr")
	}

	MessageRegistry.MapFunc(tMsg.Elem(), func(ioc *meshutil.InjectContext) interface{} {
		ev := ioc.Invoke("Event").(cellnet.Event)
		handler(ioc, ev)

		return nil
	})

}
