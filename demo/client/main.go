package main

import (
	"bufio"
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/golog"
	"os"
	"strings"
)

var log = golog.New("main")

func login() (agentAddr string) {

	loginReq, err := service.CreateConnection("login")
	if err != nil {
		log.Errorln(err)
		return
	}

	// TODO 短连接请求完毕关闭

	service.RemoteCall(loginReq, &proto.LoginREQ{
		Version:  "1.0",
		Platform: "demo",
		UID:      "1234",
	}, func(raw interface{}) {
		ack := raw.(*proto.LoginACK)

		if ack.Result == proto.ResultCode_NoError {
			agentAddr = fmt.Sprintf("%s:%d", ack.Server.IP, ack.Server.Port)
		} else {
			panic(ack.Result.String())
		}

	})

	return
}

func getAgentRequestor(agentAddr string) cellnet.Session {

	waitGameReady := make(chan cellnet.Session)
	go service.KeepConnection(agentAddr, agentAddr, waitGameReady)

	return <-waitGameReady
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

	service.RPCPairQueryFunc = proto.GetRPCPair

	svcfx.Init()

	agentAddr := login()

	if agentAddr == "" {
		return
	}

	fmt.Println("agent:", agentAddr)

	agentReq := getAgentRequestor(agentAddr)

	service.RemoteCall(agentReq, &proto.VerifyREQ{
		GameToken: "verify",
	}, func(raw interface{}) {
		ack := raw.(*proto.VerifyACK)

		fmt.Println(ack)
	})

	ReadConsole(func(s string) {

		service.RemoteCall(agentReq, &proto.ChatREQ{
			Content: s,
		}, func(raw interface{}) {
			ack := raw.(*proto.ChatACK)

			fmt.Println(ack)
		})

	})

}
