package login

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/agent/model"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
)

func init() {
	proto.Handle_Login_LoginREQ = func(ev cellnet.Event) {

		//msg := ev.Message().(*proto.LoginREQ)

		// TODO 第三方请求验证及信息拉取

		var ack proto.LoginACK

		gameList := discovery.Default.Query(model.FrontendName)
		if len(gameList) == 0 {

			ack.Result = proto.ResultCode_GameNotReady

			service.Reply(ev, &ack)
			return
		}

		// TODO 按照游戏负载选择游戏地址
		finalDesc := gameList[0]

		ack.Server.IP = finalDesc.Host
		ack.Server.Port = int32(finalDesc.Port)

		service.Reply(ev, &ack)
	}
}
