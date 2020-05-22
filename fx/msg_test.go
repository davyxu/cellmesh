package fx

import (
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"

	"testing"
)

func TestMessage(t *testing.T) {

	RegisterMessage(new(proto.LoginREQ), func(ioc *meshutil.InjectContext, ev cellnet.Event) {

		msg := ev.Message().(*proto.LoginREQ)

		t.Log(msg)

	})

	eventCallback := MakeIOCEventHandler(MessageRegistry)

	ev := &cellnet.RecvMsgEvent{
		Ses: nil,
		Msg: &proto.LoginREQ{
			Version: "1.0",
		},
	}

	eventCallback(ev)
}
