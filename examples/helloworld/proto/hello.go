package proto

import (
	"fmt"
	"github.com/davyxu/cellmicro"
	"github.com/davyxu/cellmicro/svc"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/json"
	"reflect"
)

type HelloREQ struct {
	Name string
}

type HelloACK struct {
	Message string
}

func (self *HelloREQ) String() string { return fmt.Sprintf("%+v", *self) }
func (self *HelloACK) String() string { return fmt.Sprintf("%+v", *self) }

// 客户端请求
func Hello(req *HelloREQ) (ack *HelloACK, err error) {
	ack = &HelloACK{}
	err = svc.Request(req, ack)
	return
}

// 服务器注册
func RegisterHelloHandler(s svc.Service, userHandler func(req *HelloREQ, ack *HelloACK)) {

	s.AddHandler("proto.HelloREQ", func(event *svc.Event) {

		userHandler(event.Request.(*HelloREQ), event.Response.(*HelloACK))

	})
}

func init() {

	cellmicro.RegisterRequestPair(&cellnet.MessageMeta{
		Type: reflect.TypeOf((*HelloREQ)(nil)).Elem(),
	}, &cellnet.MessageMeta{
		Type: reflect.TypeOf((*HelloACK)(nil)).Elem(),
	})
}
