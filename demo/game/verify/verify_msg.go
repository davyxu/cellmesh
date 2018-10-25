package verify

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/agent/api"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
)

func init() {

	proto.Handle_Game_VerifyREQ = api.HandleBackendMessage(func(ev cellnet.Event, cid proto.ClientID) {

		msg := ev.Message().(*proto.VerifyREQ)

		fmt.Printf("verfiy: %+v \n", msg.GameToken)

		ev.Session().Send(&proto.BindBackendACK{ID: cid.ID})

		service.Reply(ev, &proto.VerifyACK{})
	})
}
