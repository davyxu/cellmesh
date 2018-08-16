package svcfx

import (
	"flag"
	"github.com/davyxu/cellmesh/broker"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/consul"
	"github.com/davyxu/cellmesh/svcfx/model"
	"github.com/davyxu/cellnet/msglog"
	_ "github.com/davyxu/cellnet/relay" // relay消息
	"github.com/davyxu/golog"
)

var (
	flagColorLog = flag.Bool("colorlog", false, "Make log in color in *nix")

	flagIDtail = flag.String("idtail", "dev", "svcname + idtail = svcid")
)

func Init() {

	flag.Parse()

	fxmodel.IDTail = *flagIDtail

	golog.SetLevelByString("consul", "info")

	discovery.Default = consulsd.NewDiscovery()
	broker.Default = broker.NewLocalBroker()

	// 彩色日志
	if *flagColorLog {
		golog.SetColorDefine(".", msglog.LogColorDefine)
		golog.EnableColorLogger(".", true)
	}
}
