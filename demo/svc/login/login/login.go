package login

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/hub/status"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
)

func init() {
	proto.Handle_Login_LoginREQ = func(ev cellnet.Event) {

		//msg := ev.Message().(*proto.LoginREQ)
		// TODO 第三方请求验证及信息拉取

		var ack proto.LoginACK

		gameSvcID := hubstatus.SelectServiceByLowUserCount("agent", "", false)
		if gameSvcID == "" {
			ack.Result = proto.ResultCode_ServerNotFound

			service.Reply(ev, &ack)
			return
		}

		gameWAN := basefx.GetRemoteServiceWANAddress("agent", gameSvcID)

		host, port, err := util.SpliteAddress(gameWAN)
		if err != nil {
			log.Errorf("invalid address: '%s' %s", gameWAN, err.Error())

			ack.Result = proto.ResultCode_ServerAddressError

			service.Reply(ev, &ack)
			return
		}

		ack.Server = &proto.ServerInfo{
			IP:   host,
			Port: int32(port),
		}

		service.Reply(ev, &ack)
	}
}
