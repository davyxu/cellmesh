package enter

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/proto"
	agentapi "github.com/davyxu/cellmesh/svc/agent/api"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
)

func init() {
	agentapi.RegisterMessage(new(proto.LoginREQ), func(ioc *meshutil.InjectContext, ev cellnet.Event) {

		fx.Reply(ev, &proto.LoginACK{})
	})
}
