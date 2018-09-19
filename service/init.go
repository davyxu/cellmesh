package service

import (
	"flag"
	"github.com/davyxu/cellmesh/broker"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/consul"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/golog"
	"strings"
)

func Init(name string) {

	procName = name

	flag.Parse()

	matchNodes = strings.Split(*flagMatchNodes, "|")
	// 匹配节点中默认添加自己的节点
	matchNodes = append(matchNodes, GetNode())

	golog.SetLevelByString("consul", "info")

	discovery.Default = consulsd.NewDiscovery()
	broker.Default = broker.NewLocalBroker()

	// 彩色日志
	if *flagColorLog {
		golog.SetColorDefine(".", msglog.LogColorDefine)
		golog.EnableColorLogger(".", true)
	}
}
