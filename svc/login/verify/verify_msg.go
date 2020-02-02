package verify

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
)

func getAgentAddress() (code proto.ResultCode, svcID, host string, port int) {

	descList := discovery.Default.Query("frontend")

	if len(descList) == 0 {
		code = proto.ResultCode_AgentNotFound
		return
	}

	agentDesc := descList[0]

	svcID = agentDesc.ID

	wanAddr := agentDesc.GetMeta("WANAddress")

	var err error
	host, port, err = util.SpliteAddress(wanAddr)
	if err != nil {
		code = proto.ResultCode_AgentAddressError
		return
	}

	return
}

func init() {
	fx.RegisterMessage(new(proto.VerifyREQ), func(ioc *fx.InjectContext, ev cellnet.Event) {
		//msg := ev.Message().(*proto.VerifyREQ)

		var ack proto.VerifyACK

		code, svcid, host, port := getAgentAddress()

		ack.Result = code
		if code == 0 {
			ack.GameSvcID = svcid
			ack.Server.IP = host
			ack.Server.Port = int32(port)
		}

		ev.Session().Send(&ack)
	})
}
