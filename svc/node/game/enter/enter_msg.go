package enter

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/fx/db"
	agentapi "github.com/davyxu/cellmesh/svc/node/agent/api"
	"github.com/davyxu/cellmesh/svc/proto"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"github.com/gomodule/redigo/redis"
	"reflect"
)

func loadModule(conn redis.Conn, moduleName string, model interface{}) {

	vModel := reflect.TypeOf(model)
	if vModel.Kind() != reflect.Ptr {
		panic("require model ptr")
	} else {
		vModel = vModel.Elem()
	}

	for i := 0; i < vModel.NumField(); i++ {
		modelField := vModel.Field(i)
		conn.Send("HGET", moduleName, modelField.Name)
	}

	conn.Flush()

	//for i := 0; i < vModel.NumField(); i++ {
	//	modelField := vModel.Field(i)
	//	conn.Receive()
	//}
}

func loadAccount(account string) {
	db.Redis.Operate(func(conn redis.Conn) {

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

		loadAccount(msg.Token)

		fx.Reply(ev, &proto.LoginACK{})
	})
}
