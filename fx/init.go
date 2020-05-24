package fx

import (
	"bufio"
	"fmt"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/ulog"
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

	LocalSvcID = MakeSvcID(ProcName)

	initLogger()

	if *flagUseConsole {

		go ReadConsole(func(cmd string) {
			OnCommand.Invoke(cmd)
		})

	}
}

func initLogger() {
	// 设置文件日志
	if *flagLogFile != "" {

		var maxFileSize int
		if *flagLogFileSize == "" {
			ulog.Infof("LogFile: %s", *flagLogFile)
		} else {
			var err error
			maxFileSize, err = meshutil.ParseSizeString(*flagLogFileSize)
			if err == nil {
				ulog.Infof("LogFile: %s Size: %s", *flagLogFile, *flagLogFileSize)
			} else {
				ulog.Errorf("log file size err: %s", err)
			}

		}

		ulog.Global().SetOutput(ulog.NewAsyncOutput(ulog.NewRollingOutput(*flagLogFile, maxFileSize)))
	}

	switch *flagLogFormat {
	case "json":
		ulog.Global().SetFormatter(&ulog.JSONFormatter{})
	case "text":
		textFormatter := &ulog.TextFormatter{
			EnableColor: *flagLogColor,
		}

		ulog.Global().SetReportCaller(true)

		// 彩色日志
		if *flagLogColor {
			textFormatter.ParseColorRule(msglog.LogColorDefine)
		}

		ulog.Global().SetFormatter(textFormatter)
	}

	// 设置日志级别
	if *flagLogLevel != "" {
		if lv, ok := ulog.ParseLevelString(*flagLogLevel); ok {
			ulog.SetLevel(lv)
		} else {
			ulog.Warnf("invalid log level: '%s'", *flagLogLevel)
		}
	}

	// 设置消息日志模式
	if *flagMsglogMode != "" {
		msglog.SetCurrMsgLogMode(*flagMsglogMode)
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
	ulog.Infof("Execuable: %s", os.Args[0])
	ulog.Infof("WorkDir: %s", workdir)
	ulog.Infof("ProcName: '%s'", ProcName)
	ulog.Infof("PID: %d", os.Getpid())
	ulog.Infof("Discovery: '%s'", DiscoveryAddress)
	ulog.Infof("LANIP: '%s'", util.GetLocalIP())
	ulog.Infof("WANIP: '%s'", WANIP)
	ulog.Infof("SvcGroup: '%s'", SvcGroup)
	ulog.Infof("SvcIndex: %d", SvcIndex)
}

func WaitExit() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch
}

func ReadConsole(onCmd func(cmd string)) {

	for {

		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}

		// 去掉读入内容的空白符
		text = strings.TrimSpace(text)

		if len(text) == 0 {
			continue
		}

		onCmd(text)

	}
}
