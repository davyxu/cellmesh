package backend

import (
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
)

func init() {
	proto.Handler_RouterBindUserREQ = func(event service.Event, req *proto.RouterBindUserREQ) {
		model.BindClientToBackend(event.Session(), req.ID)
	}
}
