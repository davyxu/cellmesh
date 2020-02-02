package cellmesh

import (
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/svc/memsd/api"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/golog"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func Init(name string) {

	ProcName = name

	CommandLine.Parse(os.Args[1:])

	// 开发期优先从LocalFlag作用flag
	meshutil.ApplyFlagFromFile(CommandLine, *flagFlagFile)

	CommandLine.Parse(os.Args[1:])

	// 命令行可以设置GroupName, 否则初始化为IP串
	if SvcGroup == "" {
		initGroupName()
	}

	// 命令行可以设置GroupName, 否则初始化为进程ID
	if SvcIndex == 0 {
		SvcIndex = os.Getpid()
	}

	Queue = cellnet.NewEventQueue()

	Queue.StartLoop()

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

// ip+PID的16进制数值字符串，每次启动变化
func initGroupName() {

	// 兼容ipv6
	ipParts := net.ParseIP(util.GetLocalIP())

	var sb strings.Builder
	for _, p := range ipParts {
		if p == 0 || p == 255 {
			continue
		}

		sb.WriteString(fmt.Sprintf("%d", p))
	}

	SvcGroup = sb.String()
}

func LogParameter() {
	workdir, _ := os.Getwd()
	log.Infof("Execuable: %s", os.Args[0])
	log.Infof("WorkDir: %s", workdir)
	log.Infof("ProcName: '%s'", ProcName)
	log.Infof("PID: %d", os.Getpid())
	log.Infof("Discovery: '%s'", DiscoveryAddress)
	log.Infof("LANIP: '%s'", util.GetLocalIP())
	log.Infof("WANIP: '%s'", WANIP)
	log.Infof("SvcGroup: '%s'", SvcGroup)
	log.Infof("SvcIndex: %d", SvcIndex)
}

// 连接到服务发现, 建议在service.Init后, 以及服务器逻辑开始前调用
func ConnectDiscovery() {
	log.Debugf("Connecting to discovery '%s' ...", DiscoveryAddress)
	sdConfig := memsd.DefaultConfig()
	sdConfig.Address = DiscoveryAddress
	discovery.Default = memsd.NewDiscovery()
	discovery.Default.Start(sdConfig)
}

func WaitExitSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch
}
