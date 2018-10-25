package service

import (
	"errors"
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

var (
	ErrInvalidRelayMessage         = errors.New("invalid relay message")
	ErrInvalidRelayPassthroughType = errors.New("invalid relay passthrough type")
)

// 获取Event中relay的透传数据
func GetPassThrough(ev cellnet.Event, ptrList ...interface{}) error {
	if relayEvent, ok := ev.(*relay.RecvMsgEvent); ok {

		for _, ptr := range ptrList {

			switch valuePtr := ptr.(type) {
			case *int64:
				*valuePtr = relayEvent.PassThroughAsInt64()
			case *[]int64:
				*valuePtr = relayEvent.PassThroughAsInt64Slice()
			case *string:
				*valuePtr = relayEvent.PassThroughAsString()
			default:
				return ErrInvalidRelayPassthroughType
			}
		}

		return nil
	} else {
		return ErrInvalidRelayMessage
	}

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
