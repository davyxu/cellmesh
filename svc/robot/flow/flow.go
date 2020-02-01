package flow

import (
	"github.com/davyxu/cellmesh/svc/robot/model"
	"github.com/davyxu/golexer"
)

func BackgroundProc(r *model.Robot, msg interface{}) bool {

	// 异步收取全局的封包, 例如model同步等
	//switch ack := msg.(type) {
	//}

	return false
}

func Main(r *model.Robot) {
	defer golexer.ErrorCatcher(func(e error) {
		log.Errorln(e)
	})

	// 模拟异步全局收消息处理
	r.SetBackgroundRecv(func(msg interface{}) bool {

		return BackgroundProc(r, msg)
	})
}
