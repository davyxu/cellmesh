package memsd

import (
	"errors"
	"github.com/davyxu/cellnet"
	"reflect"
	"sync"
	"time"
)

var (
	callByType sync.Map // map[reflect.Type]func(interface{})
)

type typeRPCHooker struct {
}

func (typeRPCHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	outputEvent, _, err := resolveInboundEvent(inputEvent)
	if err != nil {
		log.Errorln("rpc.resolveInboundEvent", err)
		return
	}

	return outputEvent
}

func (typeRPCHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}

func resolveInboundEvent(inputEvent cellnet.Event) (ouputEvent cellnet.Event, handled bool, err error) {
	incomingMsgType := reflect.TypeOf(inputEvent.Message()).Elem()

	if rawFeedback, ok := callByType.Load(incomingMsgType); ok {
		feedBack := rawFeedback.(chan interface{})
		feedBack <- inputEvent.Message()
		return inputEvent, true, nil
	}

	return inputEvent, false, nil
}

// callback =func(ack *YouMsgACK)
func (self *memDiscovery) remoteCall(req interface{}, callback interface{}) error {

	funcType := reflect.TypeOf(callback)
	if funcType.Kind() != reflect.Func {
		panic("callback require 'func'")
	}

	self.sesGuard.RLock()
	ses := self.ses
	self.sesGuard.RUnlock()

	if ses == nil {
		return errors.New("memsd not connected")
	}

	feedBack := make(chan interface{})

	// 获取回调第一个参数

	if funcType.NumIn() != 1 {
		panic("callback func param format like 'func(ack *YouMsgACK)'")
	}

	ackType := funcType.In(0)
	if ackType.Kind() != reflect.Ptr {
		panic("callback func param format like 'func(ack *YouMsgACK)'")
	}

	ackType = ackType.Elem()

	callByType.Store(ackType, feedBack)

	defer callByType.Delete(ackType)

	ses.Send(req)

	select {
	case ack := <-feedBack:

		vCall := reflect.ValueOf(callback)

		vCall.Call([]reflect.Value{reflect.ValueOf(ack)})
		return nil
	case <-time.After(self.config.RequestTimeout):

		return errors.New("Request time out")
	}

	return nil
}
