package verify

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
)

func init() {
	proto.Handler_VerifyREQ = func(ev service.Event, req *proto.VerifyREQ) {

		fmt.Printf("verfiy: %+v \n", req.GameToken)

		ev.Session().Send(proto.RouterBindUserREQ{ID: ev.GetContextID()})

		ev.Reply(proto.VerifyACK{})
	}
}
