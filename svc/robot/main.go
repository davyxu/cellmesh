package main

import (
	_ "github.com/davyxu/cellnet/codec/protoplus"
	_ "github.com/davyxu/cellnet/peer/tcp"
)

import (
	"flag"
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/svc/robot/flow"
	"github.com/davyxu/cellmesh/svc/robot/model"
	"github.com/davyxu/cellmesh/svc/robot/rbase"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/ulog"
	"strconv"
)

func initLogger() {
	textFormatter := &ulog.TextFormatter{
		EnableColor: *rbase.FlagShowMsgLog,
	}

	if *rbase.FlagShowMsgLog {
		textFormatter.ParseColorRule(msglog.LogColorDefine)
	}

	// 彩色日志
	ulog.Global().SetFormatter(textFormatter)

	ulog.SetLevel(ulog.DebugLevel)

	msglog.SetMsgLogRule("gamedef.PingACK", "black")

	if *rbase.FlagShowMsgLog {
		msglog.SetCurrMsgLogMode("black")
	} else {
		msglog.SetCurrMsgLogMode("mute")
	}
}

func main() {

	flag.Parse()

	initLogger()

	baseID := model.GenBaseID()

	rbase.FastExec = *rbase.FlagFastFastExec

	for i := 0; i < *rbase.FlagCount; i++ {

		r := model.NewRobot(baseID + strconv.Itoa(i))
		model.AddRobot(r)

		r.SetRecvTimeoutSec(*rbase.FlagRecvTimeOut)

		r.RunFlow(flow.Main)
	}

	fx.WaitExit()
}
