package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/memsd/api"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/golog"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func Init(name string) {

	procName = name

	CommandLine.Parse(os.Args[1:])

	// 开发期优先从LocalFlag作用flag
	meshutil.ApplyFlagFromFile(CommandLine, *flagFlagFile)

	CommandLine.Parse(os.Args[1:])

	LinkRules = ParseMatchRule(getLinkRule())

	// 设置文件日志
	if *flagLogFile != "" {

		if *flagLogFileSize == "" {
			log.Infof("LogFile: %s", *flagLogFile)
			golog.SetOutputToFile(*flagLogFile)

		} else {

			size, err := meshutil.ParseSizeString(*flagLogFileSize)
			if err == nil {
				log.Infof("LogFile: %s Size: %s", *flagLogFile, *flagLogFileSize)
				golog.SetOutputToFile(*flagLogFile, golog.OutputFileOption{
					MaxFileSize: size,
				})
			} else {
				log.Errorf("log file size err: %s", err)
			}

		}
	}

	// 彩色日志
	if *flagLogColor {
		golog.SetColorDefine(".", msglog.LogColorDefine)
		golog.EnableColorLogger(".", true)
	}

	// 设置日志级别
	if *flagLogLevel != "" {

		if rawstr := strings.Split(*flagLogLevel, "|"); len(rawstr) == 2 {

			if err := golog.SetLevelByString(rawstr[0], rawstr[1]); err != nil {
				log.Warnln("SetLevelByString:", err)
			} else {
				log.Infoln("SetLevelByString:", rawstr[0], rawstr[1])
			}
		} else {
			log.Errorln("Invalid log level cli fomat, require 'name level'")
		}
	}

	// 禁用指定消息名的消息日志
	if *flagMuteMsgLog != "" {

		if err := msglog.SetMsgLogRule(*flagMuteMsgLog, msglog.MsgLogRule_BlackList); err != nil {
			log.Errorln("SetMsgLogRule: ", err)
		} else {
			log.Infoln("SetMsgLogRule:", *flagMuteMsgLog)
		}
	}

}

func getLinkRule() string {
	if *flagLinkRule == "" {
		return *flagSvcGroup
	} else {
		return *flagLinkRule
	}
}

func LogParameter() {
	workdir, _ := os.Getwd()
	log.Infof("Execuable: %s", os.Args[0])
	log.Infof("WorkDir: %s", workdir)
	log.Infof("ProcName: '%s'", GetProcName())
	log.Infof("PID: %d", os.Getpid())
	log.Infof("Discovery: '%s'", *flagDiscoveryAddr)
	log.Infof("LinkRule: '%s'", getLinkRule())
	log.Infof("SvcGroup: '%s'", GetSvcGroup())
	log.Infof("SvcIndex: %d", GetSvcIndex())
	log.Infof("LANIP: '%s'", util.GetLocalIP())
	log.Infof("WANIP: '%s'", GetWANIP())
}

// 连接到服务发现, 建议在service.Init后, 以及服务器逻辑开始前调用
func ConnectDiscovery() {
	log.Debugf("Connecting to discovery '%s' ...", *flagDiscoveryAddr)
	sdConfig := memsd.DefaultConfig()
	sdConfig.Address = *flagDiscoveryAddr
	discovery.Default = memsd.NewDiscovery(sdConfig)
}

func WaitExitSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch
}
