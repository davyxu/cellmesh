package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/examples/helloworld/proto"
	_ "github.com/davyxu/cellmesh/svc/cellsvc"
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
