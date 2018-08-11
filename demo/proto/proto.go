package proto

import (
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/json"
	"github.com/davyxu/cellnet/util"
	"reflect"
)

// 客户端请求
func Verify(targetProvider interface{}, req *VerifyREQ, callback func(ack *VerifyACK)) error {

	return service.Request(targetProvider, req, reflect.TypeOf((*VerifyACK)(nil)).Elem(), func(response interface{}) {

		callback(response.(*VerifyACK))
	})
}

// 服务器注册
func Register_Verify(s service.Service, userHandler func(req *VerifyREQ, ack *VerifyACK)) {

	s.AddMethod("proto.VerifyREQ", &service.MethodInfo{
		RequestType: reflect.TypeOf((*VerifyREQ)(nil)).Elem(),

		NewResponse: func() interface{} {
			return &VerifyACK{}
		},
		Handler: func(event *service.Event) {

			userHandler(event.Request.(*VerifyREQ), event.Response.(*VerifyACK))

		},
	})
}

func init() {

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*VerifyREQ)(nil)).Elem(),
		ID:    int(util.StringHash("proto.VerifyREQ")),
	}).SetContext("service", "demo.game")

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*VerifyACK)(nil)).Elem(),
		ID:    int(util.StringHash("proto.VerifyACK")),
	}).SetContext("service", "demo.game")

}
