package router

import (
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/cellnet/relay"
	"reflect"
)

type RelayUpMsgHooker struct {
}

// 根据消息（方法）名，查到对应的服务
func (RelayUpMsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch incomingMsg := inputEvent.Message().(type) {
	case *cellnet.SessionAccepted:
	case *cellnet.SessionClosed:
	default:
		msgType := reflect.TypeOf(incomingMsg).Elem()

		if serviceName, ok := QuerySerivceByMsgType(msgType); ok {

			addr, err := service.QueryServiceAddress(serviceName)
			if err != nil {

				log.Warnln("Get relay service address failed ", err)
				return
			}

			ses := service.GetSession(addr)

			relay.Relay(ses, incomingMsg, inputEvent.Session().ID())

		} else {
			log.Warnf("Route target not found: %s", msgType.Name())
		}
	}

	return inputEvent
}

func (RelayUpMsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}

func init() {

	transmitter := new(tcp.TCPMessageTransmitter)
	routerHooker := new(RelayUpMsgHooker)

	proc.RegisterProcessor("demo.router", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(transmitter)
		bundle.SetHooker(routerHooker)
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
