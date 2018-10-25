package main

import (
	"bufio"
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/golog"
	"os"
	"strings"
)

var log = golog.New("main")

func login() (agentAddr string) {

	log.Debugln("Create login connection...")

	loginReq := service.CreateConnection("login")

	// TODO 短连接请求完毕关闭

	service.RemoteCall(loginReq, &proto.LoginREQ{
		Version:  "1.0",
		Platform: "demo",
		UID:      "1234",
	}, func(ack *proto.LoginACK) {

		if ack.Result == proto.ResultCode_NoError {
			agentAddr = fmt.Sprintf("%s:%d", ack.Server.IP, ack.Server.Port)
		} else {
			panic(ack.Result.String())
		}

	})

	return
}

func getAgentSession(agentAddr string) (ret cellnet.Session) {

	log.Debugln("Prepare agent connection...")

	waitGameReady := make(chan struct{})
	go service.KeepConnection("agent", agentAddr, func(ses cellnet.Session) {
		ret = ses
		waitGameReady <- struct{}{}
	}, func() {
		os.Exit(0)
	})

	<-waitGameReady

	log.Debugln("Agent connection ready")

	return
}

func ReadConsole(callback func(string)) {

	for {

		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}

		// 去掉读入内容的空白符
		text = strings.TrimSpace(text)

		callback(text)

	}
}

func main() {

	service.Init("client")

	agentAddr := login()

	if agentAddr == "" {
		return
	}

	fmt.Println("agent:", agentAddr)

	agentReq := getAgentSession(agentAddr)

	service.RemoteCall(agentReq, &proto.VerifyREQ{
		GameToken: "verify",
	}, func(ack *proto.VerifyACK) {

		fmt.Println(ack)
	})

	ReadConsole(func(s string) {

		service.RemoteCall(agentReq, &proto.ChatREQ{
			Content: s,
		}, func(ack *proto.ChatACK) {

			fmt.Println(ack)
		})

	})

}
