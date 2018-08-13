package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
)

func login() (agentAddr service.AddressSource) {

	loginReq, err := service.CreateConnection("demo.login", service.NewMsgRequestor)
	if err != nil {
		fmt.Println(err)
		return
	}

	proto.Login(loginReq, &proto.LoginREQ{
		Version:  "1.0",
		Platform: "demo",
		UID:      "1234",
	}, func(ack *proto.LoginACK) {

		agentAddr = &ack.Server
	})

	loginReq.Stop()

	return
}

func main() {

	svcfx.Init()

	agentAddr := login()

	addr := fmt.Sprintf("%s:%d", agentAddr.GetIP(), agentAddr.GetPort())

	fmt.Println("agent:", addr)

	waitGameReady := make(chan service.Requestor)
	go service.KeepConnection(service.NewMsgRequestor, addr, waitGameReady)
	<-waitGameReady

	proto.Verify(agentAddr, &proto.VerifyREQ{
		GameToken: "hello",
	}, func(ack *proto.VerifyACK) {

		fmt.Println(ack)
	})

}
