package service

import (
	"flag"
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
	log.Infof("Execuable: %s", os.Args[0])
	log.Infof("WorkDir: %s", workdir)
	log.Infof("ProcName: '%s'", GetProcName())
	log.Infof("PID: %d", os.Getpid())
	log.Infof("Discovery: '%s'", *flagDiscoveryAddr)
	log.Infof("LinkRule: '%s'", *flagLinkRule)
	log.Infof("SvcGroup: '%s'", GetSvcGroup())
	log.Infof("SvcIndex: %d", GetSvcIndex())

	if !*flagDebugMode {
		golog.SetLevelByString("consul", "info")
	}

	sdConfig := consulsd.DefaultConfig()
	sdConfig.Address = *flagDiscoveryAddr
	discovery.Default = consulsd.NewDiscovery(sdConfig)

	// 彩色日志
	if *flagColorLog {
		golog.SetColorDefine(".", msglog.LogColorDefine)
		golog.EnableColorLogger(".", true)
	}
}
