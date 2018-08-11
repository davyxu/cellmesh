package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/login/login"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
)

func main() {

	svcfx.Init()

	s := service.NewService("demo.login")

	proto.Register_Login(s, login.Login)

	err := s.Run()

	if err != nil {
		fmt.Println(err)
	}
}
