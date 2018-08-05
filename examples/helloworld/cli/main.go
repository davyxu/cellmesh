package main

import (
	"fmt"
	_ "github.com/davyxu/cellmesh/discovery/consul"
	"github.com/davyxu/cellmesh/examples/helloworld/proto"
	"github.com/davyxu/cellmesh/service"
	_ "github.com/davyxu/cellmesh/service/cell"
)

func main() {

	service.PrepareConnection("cellmicro.greating")

	ack, err := proto.Hello(&proto.HelloREQ{
		Name: "davy",
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ack)
	}

}
