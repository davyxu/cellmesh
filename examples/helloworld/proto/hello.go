package proto

import (
	"fmt"
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/service"
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

	err = service.Request("cellmicro.greating", req, reflect.TypeOf((*HelloACK)(nil)).Elem(), func(response interface{}) {

		ack = response.(*HelloACK)
	})

	return
}

// 服务器注册
func RegisterHello(s service.Service, userHandler func(req *HelloREQ, ack *HelloACK)) {

	s.AddMethod("proto.HelloREQ", &service.MethodInfo{
		RequestType: reflect.TypeOf((*HelloREQ)(nil)).Elem(),

		NewResponse: func() interface{} {
			return &HelloACK{}
		},
		Handler: func(event *service.Event) {

			userHandler(event.Request.(*HelloREQ), event.Response.(*HelloACK))

		},
	})
}

func init() {

	// 底层发送还是需要依赖cellnet
	cellmesh.RegisterRequestPair(&cellnet.MessageMeta{
		Type: reflect.TypeOf((*HelloREQ)(nil)).Elem(),
	}, &cellnet.MessageMeta{
		Type: reflect.TypeOf((*HelloACK)(nil)).Elem(),
	})
}
