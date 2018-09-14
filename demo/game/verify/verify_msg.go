package verify

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/agent/api"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
)

func init() {

	proto.Handle_Game_VerifyREQ = api.HandleBackendMessage(func(ev service.Event, cid proto.ClientID) {

		msg := ev.Message().(*proto.VerifyREQ)

		fmt.Printf("verfiy: %+v \n", msg.GameToken)

		ev.Session().Send(proto.BindBackendACK{ID: cid.ID})

		ev.Reply(proto.VerifyACK{})
	})
}
