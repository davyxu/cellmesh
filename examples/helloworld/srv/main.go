package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/examples/helloworld/proto"
	"github.com/davyxu/cellmesh/service"
	_ "github.com/davyxu/cellmesh/service/cell"
)

func Hello(req *proto.HelloREQ, ack *proto.HelloACK) {

	fmt.Printf("requst: %+v \n", req)

	ack.Message = "hello " + req.Name
}

func main() {

	s := service.NewService("cellmicro.greating")

	proto.RegisterHello(s, Hello)

	err := s.Run()

	if err != nil {
		fmt.Println(err)
	}

}
