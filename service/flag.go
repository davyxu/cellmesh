package service

import "flag"

var (
	flagColorLog = flag.Bool("colorlog", false, "Make log in color in *nix")

	flagLinkRule = flag.String("linkrule", "", "discovery other node, format like: 'svcname:tgtnode|defaultnode'")

	flagNode = flag.String("node", "dev", "node name, svcname@node = unique svcid")

	flagWANIP = flag.String("wanip", "", "client connect from extern ip")
)
