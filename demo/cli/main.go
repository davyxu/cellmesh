package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	_ "github.com/davyxu/cellmesh/discovery/consul"
	"github.com/davyxu/cellmesh/service"
	_ "github.com/davyxu/cellmesh/service/cell"
	_ "github.com/davyxu/cellnet/peer/tcp"
)

func main() {

	service.PrepareConnection("demo.agent")

	ack, err := proto.Verify("demo.agent", &proto.VerifyREQ{
		Token: "hello",
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ack)
	}
}
