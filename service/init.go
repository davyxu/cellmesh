package service

import (
	"flag"
	"github.com/davyxu/cellmesh/broker"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/consul"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/golog"
	"os"
)

func Init(name string) {

	procName = name

	flag.Parse()

	if *flagLinkRule == "" {
		LinkRules = ParseMatchRule("dev") // 默认匹配dev组
	} else {
		LinkRules = ParseMatchRule(*flagLinkRule)
	}

	workdir, _ := os.Getwd()
	log.Infoln("cellmesh initializing....")
	log.Infof("ProcName: '%s'", GetProcName())
	log.Infof("SvcIndex: %d", GetSvcIndex())
	log.Infof("SvcGroup: '%s'", GetSvcGroup())
	log.Infof("LinkRule: '%s'", *flagLinkRule)
	log.Infof("Execuable: %s", os.Args[0])
	log.Infof("WorkDir: %s", workdir)
	log.Infof("PID: %d", os.Getpid())

	golog.SetLevelByString("consul", "info")

	discovery.Default = consulsd.NewDiscovery()
	broker.Default = broker.NewLocalBroker()

	// 彩色日志
	if *flagColorLog {
		golog.SetColorDefine(".", msglog.LogColorDefine)
		golog.EnableColorLogger(".", true)
	}
}
