package main

import (
	"fmt"
	_ "github.com/davyxu/cellmesh/endpoint/cellep"
	"github.com/davyxu/cellmesh/examples/helloworld/proto"
)

func main() {

	ack, err := proto.Hello(&proto.HelloREQ{
		Name: "davy",
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ack)
	}

}
