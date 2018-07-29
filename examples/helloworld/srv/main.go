package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/examples/helloworld/proto"
	"github.com/davyxu/cellmesh/svc/cellsvc"
)

func Hello(req *proto.HelloREQ, ack *proto.HelloACK) {

	fmt.Printf("requst: %+v \n", req)

	ack.Message = "hello " + req.Name
}

func main() {

	s := cellsvc.NewService("cellmicro.greating")

	proto.RegisterHelloHandler(s, Hello)

	err := s.Run()

	if err != nil {
		fmt.Println(err)
	}

}
