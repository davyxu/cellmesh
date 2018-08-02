package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/endpoint/cellep"
	"github.com/davyxu/cellmesh/examples/helloworld/proto"
)

func Hello(req *proto.HelloREQ, ack *proto.HelloACK) {

	fmt.Printf("requst: %+v \n", req)

	ack.Message = "hello " + req.Name
}

func main() {

	s := cellep.NewService("cellmicro.greating")

	proto.RegisterHello(s, Hello)

	err := s.Run()

	if err != nil {
		fmt.Println(err)
	}

}
