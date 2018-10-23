package service

import (
	"flag"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/consul"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/golog"
	"os"
)

func Init(name string) {

	procName = name

	flag.Parse()

	var linkRule string

	if *flagLinkRule == "" {
		linkRule = *flagSvcGroup
	} else {
		linkRule = *flagLinkRule
	}

	workdir, _ := os.Getwd()
	log.Infof("Execuable: %s", os.Args[0])
	log.Infof("WorkDir: %s", workdir)
	log.Infof("ProcName: '%s'", GetProcName())
	log.Infof("PID: %d", os.Getpid())
	log.Infof("Discovery: '%s'", *flagDiscoveryAddr)
	log.Infof("LinkRule: '%s'", linkRule)
	log.Infof("SvcGroup: '%s'", GetSvcGroup())
	log.Infof("SvcIndex: %d", GetSvcIndex())
	log.Infof("LANIP: '%s'", util.GetLocalIP())
	log.Infof("WANIP: '%s'", GetWANIP())

	LinkRules = ParseMatchRule(linkRule)

	log.Debugln("Connect to consul...")
	sdConfig := consulsd.DefaultConfig()
	sdConfig.Address = *flagDiscoveryAddr
	discovery.Default = consulsd.NewDiscovery(sdConfig)

	// 彩色日志
	if *flagColorLog {
		golog.SetColorDefine(".", msglog.LogColorDefine)
		golog.EnableColorLogger(".", true)
	}
}
