package login

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery"
)

func Login(req *proto.LoginREQ, ack *proto.LoginACK) {

	// TODO DB请求验证登录

	gameList, err := discovery.Default.Query("demo.game")
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
