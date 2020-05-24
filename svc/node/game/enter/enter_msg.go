package enter

import (
	"github.com/davyxu/cellmesh/fx"
	agentapi "github.com/davyxu/cellmesh/svc/node/agent/api"
	"github.com/davyxu/cellmesh/svc/proto"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
)

func loadAccount(account string) {

}

func init() {
	agentapi.RegisterMessage(new(proto.AgentBindBackendREQ), func(ioc *meshutil.InjectContext, ev cellnet.Event) {
		msg := ev.Message().(*proto.AgentBindBackendREQ)
		fx.Reply(ev, &proto.AgentBindBackendACK{
			NodeID: msg.NodeID,
		})
	})

	agentapi.RegisterMessage(new(proto.LoginREQ), func(ioc *meshutil.InjectContext, ev cellnet.Event) {

		msg := ev.Message().(*proto.LoginREQ)

		loadAccount(msg.Token)

		fx.Reply(ev, &proto.LoginACK{})
	})
}
