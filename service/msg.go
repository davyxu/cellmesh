package service

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/relay"
	"github.com/davyxu/cellnet/util"
	"reflect"
)

type ServiceIdentifyACK struct {
	SvcName string
	SvcID   string
}

func (self *ServiceIdentifyACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*ServiceIdentifyACK)(nil)).Elem(),
		ID:    int(util.StringHash("service.ServiceIdentifyACK")),
	})
}

// 获取Event中relay的透传数据
func GetPassThrough(ev cellnet.Event) interface{} {
	if relayEvent, ok := ev.(*relay.RecvMsgEvent); ok {
		return relayEvent.PassThrough
	}

	return nil
}

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
