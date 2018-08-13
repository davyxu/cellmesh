package login

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery"
)

func Login(req *proto.LoginREQ, ack *proto.LoginACK) {

	// TODO 第三方请求验证及信息拉取

	gameList, err := discovery.Default.Query("demo.agent")
	if err != nil || len(gameList) == 0 {
		ack.Result = proto.ResultCode_GameNotReady
		return
	}

	// TODO 按照游戏负载选择游戏地址
	finalDesc := gameList[0]

	ack.Server.IP = finalDesc.Address
	ack.Server.Port = int32(finalDesc.Port)

	// 没有错误
	ack.Result = 0
}
