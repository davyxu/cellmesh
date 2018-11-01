package subscribe

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/hub/model"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/relay"
)

func init() {

	relay.SetBroadcaster(func(event *relay.RecvMsgEvent) {

		if channelName := event.PassThroughAsString(); channelName != "" {

			model.VisitSubscriber(channelName, func(ses cellnet.Session) bool {

				relay.Relay(ses, event.Message(), channelName)

				return true
			})
		}

	})

	proto.Handle_Hub_SubscribeChannelREQ = func(ev cellnet.Event) {

		msg := ev.Message().(*proto.SubscribeChannelREQ)
		model.AddSubscriber(msg.Channel, ev.Session())

		log.Infof("channel add: '%s', sesid: %d", msg.Channel, ev.Session().ID())

		ev.Session().Send(&proto.SubscribeChannelACK{
			Channel: msg.Channel,
		})
	}

	proto.Handle_Hub_Default = func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionClosed:

			model.RemoveSubscriber(ev.Session(), func(chanName string) {
				log.Infof("channel remove: '%s', sesid: %d", chanName, ev.Session().ID())
			})
		}
	}
}
