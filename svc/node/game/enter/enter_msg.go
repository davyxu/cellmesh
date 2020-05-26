package enter

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/fx/db"
	"github.com/davyxu/cellmesh/svc/actor"
	agentapi "github.com/davyxu/cellmesh/svc/node/agent/api"
	"github.com/davyxu/cellmesh/svc/proto"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"github.com/gomodule/redigo/redis"
)

func loadAccount(account string) db.ResultCode {
	return db.Redis.Operate(func(conn redis.Conn) {

		var acc actor.Account

		ser := db.NewModelList(conn)

		if ser.Load(&acc, account) == db.ErrModelNotExists {
			acc.Account = account

			ser.Save(&acc, account)
		}
	})
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

		var ack proto.LoginACK
		ack.Code = loadAccount(msg.Token)

		fx.Reply(ev, &ack)
	})
}
