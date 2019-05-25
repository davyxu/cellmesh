package basefx

import (
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
)

// 初始化框架
func Init(procName string) {

	msglog.SetCurrMsgLogMode(msglog.MsgLogMode_BlackList)
	msglog.SetMsgLogRule("proto.PingACK", msglog.MsgLogRule_BlackList)
	msglog.SetMsgLogRule("proto.SvcStatusACK", msglog.MsgLogRule_BlackList)

	fxmodel.Queue = cellnet.NewEventQueue()

	fxmodel.Queue.StartLoop()

	service.Init(procName)

	service.ConnectDiscovery()
}

// 等待退出信号
func StartLoop(onReady func()) {

	fxmodel.CheckReady()

	if onReady != nil {
		cellnet.QueuedCall(fxmodel.Queue, onReady)
	}

	service.WaitExitSignal()
}

// 退出处理
func Exit() {
	fxmodel.StopAllService()
}
