package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
)

func login() (agentAddr string) {

	loginReq, err := service.CreateConnection("demo.login", service.NewMsgRequestor)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer loginReq.Stop()

	proto.Login(loginReq, &proto.LoginREQ{
		Version:  "1.0",
		Platform: "demo",
		UID:      "1234",
	}, func(ack *proto.LoginACK) {

		agentAddr = fmt.Sprintf("%s:%d", ack.Server.IP, ack.Server.Port)
	})

	return
}

func getAgentRequestor(agentAddr string) service.Requestor {
	waitGameReady := make(chan service.Requestor)
	go service.KeepConnection(service.NewMsgRequestor, agentAddr, waitGameReady)
	return <-waitGameReady
}

func main() {

	svcfx.Init()

	agentAddr := login()

	fmt.Println("agent:", agentAddr)

	agentReq := getAgentRequestor(agentAddr)

	proto.Verify(agentReq, &proto.VerifyREQ{
		GameToken: "hello",
	}, func(ack *proto.VerifyACK) {

		fmt.Println(ack)
	})

}
