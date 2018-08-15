package service

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/relay"
	"reflect"
	"sync"
)

type Dispatcher struct {
	svcByName sync.Map // map[reflect.Type]*endpoint.MethodInfo
}

func (self *Dispatcher) AddCall(name string, svc *MethodInfo) {

	self.svcByName.Store(svc.RequestType, svc)
}

func (self *Dispatcher) Invoke(ev cellnet.Event) {

	msgType := reflect.TypeOf(ev.Message()).Elem()

	if svcRaw, ok := self.svcByName.Load(msgType); ok {

		svc := svcRaw.(*MethodInfo)

		e := &Event{
			Request: ev.Message(),
			Session: ev.Session(),
		}

		if relayEvent, ok := ev.(*relay.RecvMsgEvent); ok {
			e.ContextID = relayEvent.ContextID
		}

		if replyEv, ok := ev.(ReplyEvent); ok {

			e.Response = svc.NewResponse()

			svc.Handler(e)

			replyEv.Reply(e.Response)

		} else {

			svc.Handler(e)

		}
	}

}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}
