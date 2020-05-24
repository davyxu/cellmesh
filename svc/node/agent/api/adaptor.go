package agentapi

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/fx/link"
	"github.com/davyxu/cellmesh/svc/proto"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"reflect"
)

func RegisterMessage(msgTypeObj interface{}, handler func(ioc *meshutil.InjectContext, ev cellnet.Event)) {

	tMsg := reflect.TypeOf(msgTypeObj)
	if tMsg.Kind() != reflect.Ptr {
		panic("require msg ptr")
	}

	fx.MessageRegistry.MapFunc(tMsg.Elem(), func(ioc *meshutil.InjectContext) interface{} {
		ev := ioc.Invoke("Event").(cellnet.Event)
		//cid := ioc.Invoke("CID").(proto.AgentClientID)
		handler(ioc, ev)

		return nil
	})

}

func invokeAgentMessage(ev cellnet.Event) {
	ioc := meshutil.NewInjectContext()

	ioc.SetParent(fx.MessageRegistry)

	ioc.MapFunc("Event", func(ioc *meshutil.InjectContext) interface{} {
		return ev
	})

	ioc.MapFunc("CID", func(ioc *meshutil.InjectContext) interface{} {

		aev := ev.(*AgentMsgEvent)

		var cid proto.AgentClientID
		cid.SessionID = aev.ClientID

		if desc := link.DescByLink(aev.Session()); desc != nil {
			cid.NodeID = desc.ID
		}

		return cid
	})

	tMsg := reflect.TypeOf(ev.Message())
	if tMsg.Kind() == reflect.Ptr {
		tMsg = tMsg.Elem()
	}
	ioc.TryInvoke(tMsg)
}
