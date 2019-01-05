package main

import (
	"flag"
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/memsd/api"
	"os"
)

var (
	flagCmd  = flag.String("cmd", "", "sub command, empty to launch memsd service")
	flagAddr = flag.String("addr", "", "service discovery address")
)

type DiscoveryExtend interface {
	discovery.Discovery

	QueryAll() (ret []*discovery.ServiceDesc)

	ClearKey()

	ClearService()

	GetRawValueList(prefix string) (ret []discovery.ValueMeta)
}

func initSD() DiscoveryExtend {
	config := memsd.DefaultConfig()
	if *flagAddr != "" {
		config.Address = *flagAddr
	}

	return memsd.NewDiscovery(config).(DiscoveryExtend)
}

func main() {

	flag.Parse()

	switch *flagCmd {
	case "": // addr
		startSvc()
	case "viewsvc": // addr
		viewSvc()
	case "viewkey": // addr
		viewKey()
	case "clearsvc": // addr
		clearSvc()
	case "clearkey": // addr
		clearKey()
	case "getvalue":
		if flag.NArg() < 1 {
			fmt.Println("getvalue <key>")
			os.Exit(1)
		}
		getValue(flag.Arg(0))
	case "setvalue":
		if flag.NArg() < 2 {
			fmt.Println("setvalue <key> <value>")
			os.Exit(1)
		}

		setValue(flag.Arg(0), flag.Arg(1))
	default:
		fmt.Printf("Unknown command '%s'\n", *flagCmd)
		os.Exit(1)
	}
}
