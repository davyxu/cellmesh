package router

import (
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/cellnet/relay"
	"reflect"
)

type RelayUpMsgHooker struct {
}

// 从客户端接收到的消息
func (RelayUpMsgHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	switch incomingMsg := inputEvent.Message().(type) {
	case *cellnet.SessionAccepted:
	case *cellnet.SessionClosed:
	default:
		msgType := reflect.TypeOf(incomingMsg).Elem()

		if serviceName, ok := QuerySerivceByMsgType(msgType); ok {

			// TODO 会话绑定，负载均衡
			addr, err := service.QueryServiceAddress(serviceName)
			if err != nil {

				log.Warnln("Get relay service address failed ", err)
				return
			}

			ses := service.GetSession(addr)

			// 服务没有连接
			if ses == nil {
				log.Warnf("service '%s' not reachable, addr: %s", serviceName, addr)
				return nil
			}

			// 透传消息
			relay.Relay(ses, incomingMsg, inputEvent.Session().ID())

		} else {

			log.Warnf("Route target not found: %s", msgType.Name())
		}
	}

	return inputEvent
}

// 发送到客户端的消息
func (RelayUpMsgHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}

func init() {

	// 从后端服务器收到的消息
	relay.SetBroadcaster(func(event *relay.RecvMsgEvent) {

		for _, sesID := range event.ContextID {

			ses := clientListener.(peer.SessionManager).GetSession(sesID)
			if ses == nil {
				continue
			}

			ses.Send(event.Msg)
		}

	})

	transmitter := new(tcp.TCPMessageTransmitter)
	routerHooker := new(RelayUpMsgHooker)

	proc.RegisterProcessor("demo.router", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(transmitter)
		bundle.SetHooker(routerHooker)
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
