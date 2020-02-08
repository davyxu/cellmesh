package verify

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/link"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellmesh/rpc"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/ulog"
)

func getAgentAddress() (code proto.ResultCode, host string, port int) {

	descList := discovery.Global.Query("frontend")

	if len(descList) == 0 {
		code = proto.ResultCode_AgentNotFound
		return
	}

	// TODO 挑选低负载agent
	agentDesc := descList[0]
	wanAddr := agentDesc.GetMeta("WANAddress")
	var err error
	host, port, err = util.SpliteAddress(wanAddr)
	if err != nil {
		code = proto.ResultCode_AgentAddressError
		return
	}

	return
}

func getGameSvcID() (code proto.ResultCode, svcID string) {
	descList := discovery.Global.Query("game")

	if len(descList) == 0 {
		code = proto.ResultCode_GameNotFound
		return
	}

	// TODO 挑选低负载game
	svcID = descList[0].ID

	return
}

func init() {

	fx.OnLoad.Add(func(args ...interface{}) {

		rpc.New(link.OneLink("hub")).Request(&proto.TestREQ{
			Dummy: "hello",
		}).RecvWait(func(resp *rpc.Respond) {
			if resp.Error != nil {
				ulog.Errorln(resp.Error)
				return
			}

			ulog.Debugf("%+V", resp)
		})

	})

	fx.RegisterMessage(new(proto.VerifyREQ), func(ioc *fx.InjectContext, ev cellnet.Event) {
		//msg := ev.Message().(*proto.VerifyREQ)

		code, host, port := getAgentAddress()
		if code != 0 {
			fx.Reply(ev, &proto.VerifyACK{
				Code: code,
			})
			return
		}

		code, svcid := getGameSvcID()
		if code != 0 {
			fx.Reply(ev, &proto.VerifyACK{
				Code: code,
			})
			return
		}

		var ack proto.VerifyACK
		ack.Server.IP = host
		ack.Server.Port = int32(port)
		ack.SvcID = svcid
		ack.Code = code

		fx.Reply(ev, &ack)
	})
}
