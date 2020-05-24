package verify

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/fx/link"
	"github.com/davyxu/cellmesh/fx/rpc"
	"github.com/davyxu/cellmesh/svc/proto"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/ulog"
)

func getAgentAddress() (code proto.ResultCode, host string, port int) {

	descList := link.DescListByName("frontend")

	if len(descList) == 0 {
		code = proto.ResultCode_AgentNotFound
		return
	}

	// TODO 挑选低负载agent
	agentDesc := descList[0]
	wanAddr := agentDesc.GetMeta("WAN")
	var err error
	host, port, err = util.SpliteAddress(wanAddr)
	if err != nil {
		code = proto.ResultCode_AgentAddressError
		return
	}

	return
}

func getGameNodeID() (code proto.ResultCode, svcID string) {
	descList := link.DescListByName("game")

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

		rpc.New(link.LinkByName("hub")).Request(&proto.TestREQ{
			Dummy: "hello",
		}).RecvWait(func(resp *rpc.Respond) {
			if resp.Error != nil {
				ulog.Errorln(resp.Error)
				return
			}

			ulog.Debugf("%+v", resp)
		})

	})

	fx.RegisterMessage(new(proto.VerifyREQ), func(ioc *meshutil.InjectContext, ev cellnet.Event) {
		//msg := ev.Message().(*proto.VerifyREQ)

		code, host, port := getAgentAddress()
		if code != 0 {
			fx.Reply(ev, &proto.VerifyACK{
				Code: code,
			})
			return
		}

		code, nodeid := getGameNodeID()
		if code != 0 {
			fx.Reply(ev, &proto.VerifyACK{
				Code: code,
			})
			return
		}

		var ack proto.VerifyACK
		ack.Server.IP = host
		ack.Server.Port = int32(port)
		ack.NodeID = nodeid
		ack.Code = code

		fx.Reply(ev, &ack)
	})
}
