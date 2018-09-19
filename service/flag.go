package service

import "flag"

var (
	flagColorLog = flag.Bool("colorlog", false, "Make log in color in *nix")

	flagMatchNodes = flag.String("matchnodes", "", "discovery other node, split by |")

	flagNode = flag.String("node", "dev", "node name, svcname@node = unique svcid")

	flagWANIP = flag.String("wanip", "", "client connect from extern ip")
)
