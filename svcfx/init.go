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
	"strings"
)

var (
	flagColorLog = flag.Bool("colorlog", false, "Make log in color in *nix")

	flagMatchNodes = flag.String("matchnodes", "", "discovery other node, split by |")

	flagNode = flag.String("node", "dev", "node name, svcname@node = unique svcid")
)

func Init() {

	flag.Parse()

	fxmodel.MatchNodes = strings.Split(*flagMatchNodes, "|")
	fxmodel.Node = *flagNode

	// 匹配节点中默认添加自己的节点
	fxmodel.MatchNodes = append(fxmodel.MatchNodes, fxmodel.Node)

	golog.SetLevelByString("consul", "info")

	discovery.Default = consulsd.NewDiscovery()
	broker.Default = broker.NewLocalBroker()

	// 彩色日志
	if *flagColorLog {
		golog.SetColorDefine(".", msglog.LogColorDefine)
		golog.EnableColorLogger(".", true)
	}
}
