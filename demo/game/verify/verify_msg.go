package verify

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
)

func Verify(ev *service.Event, req *proto.VerifyREQ, ack *proto.VerifyACK) {

	fmt.Printf("verfiy: %+v \n", req.GameToken)

	ack.Result = 0
}
