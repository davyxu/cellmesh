package backend

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
)

func init() {
	proto.Handler_RouterBindUserREQ = func(event service.Event, req *proto.RouterBindUserREQ) {
		fmt.Println("bind user", req.Token)
	}
}
