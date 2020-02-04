package fx

import (
	"github.com/davyxu/cellnet"
	"reflect"
)

// 回复event来源一个消息
func Reply(ev cellnet.Event, msg interface{}) {

	type replyEvent interface {
		Reply(msg interface{})
	}

	if replyEv, ok := ev.(replyEvent); ok {
		replyEv.Reply(msg)
	} else {
		panic("Require 'ReplyEvent' to reply event")
	}
}

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
