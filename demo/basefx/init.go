package basefx

import (
	"github.com/davyxu/cellmesh/demo/basefx/model"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"os"
	"os/signal"
	"syscall"
)

// 初始化框架
func Init(procName string) {

	fxmodel.Queue = cellnet.NewEventQueue()

	fxmodel.Queue.StartLoop()

	service.Init(procName)
}

// 等待退出信号
func StartLoop() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch
}

// 退出处理
func Exit() {
	StopAllPeers()
}
