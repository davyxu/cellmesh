package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/login/login"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

func main() {

	svcfx.Init()

	s := service.NewService("demo.login")

	proto.Serve_Login(s, login.Login)

	err := s.Run()

	if err != nil {
		fmt.Println(err)
	}
}
