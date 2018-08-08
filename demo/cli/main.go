package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
	"reflect"
	"time"
)

func main() {

	svcfx.Init()

	addr, err := service.QueryServiceAddress("demo.agent")
	if err != nil {
		fmt.Println(err)
		return
	}

	req := service.NewMsgRequestor(addr, nil)
	req.Start()
	for !req.IsReady() {

		time.Sleep(time.Second)
		req.Stop()
	}

	err = req.Request(&proto.VerifyREQ{
		Token: "hello",
	}, reflect.TypeOf((*proto.VerifyACK)(nil)).Elem(), func(ack interface{}) {

		fmt.Println(ack)

	})

	if err != nil {
		fmt.Println(err)
	}
}
