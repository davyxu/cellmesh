package svcfx

import (
	"flag"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/consul"
	"github.com/davyxu/cellnet/msglog"
	_ "github.com/davyxu/cellnet/relay" // relay消息
	"github.com/davyxu/golog"
)

var (
	flagColorLog = flag.Bool("colorlog", false, "Make log in color in *nix")
)

func Init() {

	flag.Parse()

	discovery.Default = consulsd.NewDiscovery()

	// 彩色日志
	if *flagColorLog {
		golog.SetColorDefine(".", msglog.LogColorDefine)
		golog.EnableColorLogger(".", true)
	}
}
