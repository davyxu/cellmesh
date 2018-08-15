package main

import (
	"github.com/davyxu/cellmesh/demo/agent/backend"
	"github.com/davyxu/cellmesh/demo/agent/frontend"
	"github.com/davyxu/cellmesh/demo/proto"
	_ "github.com/davyxu/cellmesh/demo/proto" // 进入协议
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	svcfx.Init()

	// 网关主动连接后台的服务器，因为后台服务是"服务"
	go service.PrepareConnection("demo.game", service.NewRPCRequestor, func(desc *discovery.ServiceDesc, requestor service.Requestor) {

	})

	s := service.NewService("demo.router")
	proto.Serve_RouterBindUser(s, backend.RouterBindUser)
	go s.Run()

	frontend.Run()
}
