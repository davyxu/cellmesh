package flow

import (
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellmesh/svc/robot/link"
	"github.com/davyxu/cellmesh/svc/robot/model"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/ulog"
)

func BackgroundProc(r *model.Robot, msg interface{}) bool {

	// 异步收取全局的封包, 例如model同步等
	//switch ack := msg.(type) {
	//}

	return false
}

func Main(r *model.Robot) {
	defer util.ErrorCatcher(func(e error) {
		ulog.Errorln(e)
	})

	// 模拟异步全局收消息处理
	r.SetBackgroundRecv(func(msg interface{}) bool {

		return BackgroundProc(r, msg)
	})

	link.ConnectTCP(r, "login", util.GetLocalIP()+":8001")

	r.Send("login", &proto.VerifyREQ{})
	ack := r.Recv("proto.VerifyACK").(*proto.VerifyACK)
	ulog.Debugf("%+v", ack)

	link.ConnectTCP(r, "agent", util.JoinAddress(ack.Server.IP, int(ack.Server.Port)))
}
