package main

import (
	"bufio"
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/timer"
	"github.com/davyxu/golog"
	"os"
	"strings"
	"time"
)

var log = golog.New("main")

// 登录服务器是一个短连接服务器,获取到网关连接后就断开
func login() (agentAddr string) {

	log.Debugln("Create login connection...")

	// 这里为了方便,使用服务发现来连接login, 真正的客户端不应该使用服务发现,而是使用固定的登录地址连接登录服务器
	loginReq := service.CreateConnection("login")

	// TODO 短连接请求完毕关闭

	// 封装连接和接收以及超时的过程,真正的客户端可以添加这套机制
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

	msglog.BlockMessageLog("proto.PingACK")

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

	//
	timer.NewLoop(nil, time.Second*5, func(loop *timer.Loop) {
		agentReq.Send(&proto.PingACK{})

	}, nil).Start()

	ReadConsole(func(s string) {

		service.RemoteCall(agentReq, &proto.ChatREQ{
			Content: s,
		}, func(ack *proto.ChatACK) {

			fmt.Println(ack.Content)
		})

	})

}
