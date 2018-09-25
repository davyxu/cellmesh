package service

import "flag"

var (
	flagColorLog = flag.Bool("colorlog", false, "Make log in color in *nix")

	flagLinkRule = flag.String("linkrule", "", "discovery other node then connect it, format like: 'svcname:tgtnode|defaultnode'")

	flagSvcGroup = flag.String("svcgroup", "dev", "represent one group server")

	flagSvcIndex = flag.Int("svcindex", 0, "multi proc in group use index to seperate each other")

	flagWANIP = flag.String("wanip", "", "client connect from extern ip")

	flagDebugMode = flag.Bool("debug", false, "show debug info")
)
