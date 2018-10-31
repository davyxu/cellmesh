package service

import "flag"

var (

	// 服务发现地址
	flagDiscoveryAddr = flag.String("sdaddr", "127.0.0.1:8500", "Discovery address")

	// 服务发现规则
	flagLinkRule = flag.String("linkrule", "", "discovery other node then connect it, format like: 'svcname:tgtnode|defaultnode'")

	// 服务所在组
	flagSvcGroup = flag.String("svcgroup", "dev", "represent one group server")

	// 服务索引
	flagSvcIndex = flag.Int("svcindex", 0, "multi proc in group use index to seperate each other")

	// 设置外网IP
	flagWANIP = flag.String("wanip", "", "client connect from extern ip")

	// 对日志上色
	flagLogColor = flag.Bool("logcolor", false, "Make log in color in *nix")

	// 将日志输出到文件
	flagLogFile = flag.String("logfile", "", "log file name")

	// 设置日志级别，格式: 日志名称 日志级别， 名称支持正则表达式
	flagLogLevel = flag.String("loglevel", "", "Set log level, format: 'name|level', name support regexp, level can be error, info")

	// 屏蔽消息日志
	flagMuteMsgLog = flag.String("mutemsglog", "", "do not show msg log, splite by '|', support regexp")
)
