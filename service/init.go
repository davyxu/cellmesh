package service

import (
	"flag"
	"github.com/davyxu/cellmesh/broker"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/consul"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/golog"
)

func Init(name string) {

	procName = name

	flag.Parse()

	if *flagLinkRule == "" {
		LinkRules = ParseMatchRule("dev") // 默认匹配dev组
	} else {
		LinkRules = ParseMatchRule(*flagLinkRule)
	}

	golog.SetLevelByString("consul", "info")

	discovery.Default = consulsd.NewDiscovery()
	broker.Default = broker.NewLocalBroker()

	// 彩色日志
	if *flagColorLog {
		golog.SetColorDefine(".", msglog.LogColorDefine)
		golog.EnableColorLogger(".", true)
	}
}
