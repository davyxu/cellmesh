package hubapi

import (
	"github.com/davyxu/cellmesh/demo/basefx"
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/hub/model"
	"github.com/davyxu/cellnet/relay"
)

// 传入你的服务名, 连接到hub
func ConnectToHub(hubReady func()) {

	model.OnHubReady = hubReady
	basefx.CreateCommnicateConnector(fxmodel.ServiceParameter{
		SvcName:      "hub",
		NetProcName:  "tcp.hub",
		MaxConnCount: 1,
	})
}

func Subscribe(channel string) {

	if model.HubSession == nil {
		log.Errorf("hub session not ready, channel: %s", channel)
		return
	}

	model.HubSession.Send(&proto.SubscribeChannelREQ{
		Channel: channel,
	})
}

func Publish(channel string, msg interface{}) {

	if model.HubSession == nil {
		log.Errorf("hub session not ready, channel: %s", channel)
		return
	}

	relay.Relay(model.HubSession, msg, channel)
}
