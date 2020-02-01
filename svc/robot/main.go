package main

import (
	"github.com/davyxu/cellmesh"
	"github.com/davyxu/cellmesh/svc/robot/flow"
	"github.com/davyxu/cellmesh/svc/robot/model"
	"github.com/davyxu/cellmesh/svc/robot/rbase"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/golog"
	"strconv"
)

var log = golog.New("main")

func main() {

	// 异步写日志
	golog.EnableASyncWrite()

	// 精确到毫秒
	golog.VisitLogger("[.]*", func(logger *golog.Logger) bool {

		logger.SetParts(golog.LogPart_CurrLevel, golog.LogPart_Name, golog.LogPart_TimeMS)

		return true
	})

	msglog.SetCurrMsgLogMode(msglog.MsgLogMode_BlackList)
	//msglog.SetMsgLogRule("gamedef.PingACK", msglog.MsgLogRule_BlackList)

	if !*rbase.FlagShowMsgLog {
		msglog.SetCurrMsgLogMode(msglog.MsgLogMode_Mute)
	}

	baseID := model.GenBaseID()

	rbase.FastExec = *rbase.FlagFastFastExec

	for i := 0; i < *rbase.FlagCount; i++ {

		r := model.NewRobot(baseID + strconv.Itoa(i))
		model.AddRobot(r)

		r.SetRecvTimeoutSec(*rbase.FlagRecvTimeOut)

		r.RunFlow(flow.Main)

	}

	cellmesh.WaitExitSignal()
}
