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

		agentSvcID := hubstatus.SelectServiceByLowUserCount("agent", "", false)
		if agentSvcID == "" {
			ack.Result = proto.ResultCode_AgentNotFound

			service.Reply(ev, &ack)
			return
		}

		agentWAN := basefx.GetRemoteServiceWANAddress("agent", agentSvcID)


		host, port, err := util.SpliteAddress(agentWAN)
		if err != nil {
			log.Errorf("invalid address: '%s' %s", agentWAN, err.Error())

			ack.Result = proto.ResultCode_AgentAddressError

			service.Reply(ev, &ack)
			return
		}

		ack.Server = &proto.ServerInfo{
			IP:   host,
			Port: int32(port),
		}

		ack.GameSvcID = hubstatus.SelectServiceByLowUserCount("game", "", false)

		if ack.GameSvcID == "" {
			ack.Result = proto.ResultCode_GameNotFound

			service.Reply(ev, &ack)
			return
		}

		service.Reply(ev, &ack)
	}
}
