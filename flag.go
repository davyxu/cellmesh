package cellmesh

import (
	"flag"
	"os"
)

var (
	// 独立出来避免污染工具类的flagset
	CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// 服务发现地址
	flagDiscoveryAddr = CommandLine.String("sdaddr", "127.0.0.1:8900", "Discovery address")

	// 设置外网IP
	flagWANIP = CommandLine.String("wanip", "", "Client connect from extern ip")

	// 对日志上色
	flagLogColor = CommandLine.Bool("logcolor", false, "Make log in color in *nix")

	// 将日志输出到文件
	flagLogFile = CommandLine.String("logfile", "", "Print log to file by file name")

	// 单个日志文件大小, 超过设定值时创建新文件, 0表示单文件
	flagLogFileSize = CommandLine.String("logfilesize", "", "log max file size, can use B M G to represent size")

	// 设置日志级别，格式: 日志名称 日志级别， 名称支持正则表达式
	flagLogLevel = CommandLine.String("loglevel", "", "Set log level, format: 'name|level', name support regexp, level can be error, info")

	// 屏蔽消息日志
	flagMuteMsgLog = CommandLine.String("mutemsglog", "", "Do not show msg log, splite by '|', support regexp")

	// 批量设置flag
	flagFlagFile = CommandLine.String("flagfile", "../cfg/LocalFlag.cfg", "Flagfile to init flag values")
)
