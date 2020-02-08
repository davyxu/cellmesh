package relay

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
)

func init() {

	fx.RegisterMessage(new(proto.TestREQ), func(ioc *fx.InjectContext, ev cellnet.Event) {
		msg := ev.Message().(*proto.TestREQ)

		fx.Reply(ev, &proto.TestACK{
			Dummy: msg.Dummy,
		})
	})
}
