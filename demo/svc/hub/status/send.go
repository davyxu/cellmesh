package hubstatus

import (
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/demo/svc/hub/api"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet/timer"
	"time"
)

func StartSendStatus(channelName string, updateInterval time.Duration, statusGetter func() int) {

	timer.NewLoop(fxmodel.Queue, updateInterval, func(loop *timer.Loop) {

		var ack proto.SvcStatusACK
		ack.SvcID = service.GetLocalSvcID()
		ack.UserCount = int32(statusGetter())

		hubapi.Publish(channelName, &ack)

	}, nil).Notify().Start()
}
