package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
	"time"
)

func SafeGetServiceAddress(serviceName string) service.Requestor {

RetryDiscovery:
	addr, err := service.QueryServiceAddress(serviceName)
	if err != nil {
		fmt.Println(err)

		time.Sleep(time.Second * 3)
		goto RetryDiscovery
	}

	requestor := service.NewMsgRequestor(addr, nil)
	requestor.Start()
	for !requestor.IsReady() {

		time.Sleep(time.Second)
		requestor.Stop()
	}

	return requestor
}

func login() (agentAddr service.AddressSource) {
	loginReq := SafeGetServiceAddress("demo.login")

	proto.Login(loginReq, &proto.LoginREQ{
		Version:  "1.0",
		Platform: "demo",
		UID:      "1234",
	}, func(ack *proto.LoginACK) {

		agentAddr = &ack.Server

		addr := fmt.Sprintf("%s:%d", agentAddr.GetIP(), agentAddr.GetPort())

		requestor := service.NewMsgRequestor(addr, nil)
		requestor.Start()
		for !requestor.IsReady() {

			time.Sleep(time.Second)
			requestor.Stop()
		}
	})

	return
}

func main() {

	svcfx.Init()

	agentAddr := login()

	proto.Verify(agentAddr, &proto.VerifyREQ{
		GameToken: "hello",
	}, func(ack *proto.VerifyACK) {

		fmt.Println(ack)
	})

}
