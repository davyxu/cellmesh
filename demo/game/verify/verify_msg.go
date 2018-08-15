package verify

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
)

func Verify(ev *service.Event, req *proto.VerifyREQ, ack *proto.VerifyACK) {

	fmt.Printf("verfiy: %+v \n", req.GameToken)

	ev.Session.Send(proto.RouterBindUserREQ{Token: 5115})

	ack.Result = 0
}
