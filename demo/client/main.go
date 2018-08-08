package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
	"time"
)

func main() {

	svcfx.Init()

	addr, err := service.QueryServiceAddress("demo.agent")
	if err != nil {
		fmt.Println(err)
		return
	}

	requestor := service.NewMsgRequestor(addr, nil)
	requestor.Start()
	for !requestor.IsReady() {

		time.Sleep(time.Second)
		requestor.Stop()
	}

	err = proto.Verify(requestor, &proto.VerifyREQ{
		Token: "hello",
	}, func(ack *proto.VerifyACK) {

		fmt.Println(ack)
	})

	if err != nil {
		fmt.Println(err)
	}
}
