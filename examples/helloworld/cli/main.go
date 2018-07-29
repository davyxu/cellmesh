package main

import (
	"fmt"
	"github.com/davyxu/cellmicro/examples/helloworld/proto"
	_ "github.com/davyxu/cellmicro/svc/cellsvc"
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
