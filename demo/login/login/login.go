package login

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
)

func init() {
	proto.Handle_Login_LoginREQ = func(ev service.Event) {

		//msg := ev.Message().(*proto.LoginREQ)

		// TODO 第三方请求验证及信息拉取

		var ack proto.LoginACK

		gameList, err := discovery.Default.Query("agent")
		if err != nil || len(gameList) == 0 {

			ack.Result = proto.ResultCode_GameNotReady

			ev.Reply(&ack)
			return
		}

		// TODO 按照游戏负载选择游戏地址
		finalDesc := gameList[0]

		ack.Server.IP = finalDesc.Host
		ack.Server.Port = int32(finalDesc.Port)

		ev.Reply(&ack)
	}
}
