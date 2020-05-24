package main

import (
	"github.com/davyxu/cellmesh/fx/redsd"
	_ "github.com/davyxu/cellnet/codec/protoplus"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/protoplus/msgidutil"
)

import (
	"flag"
	"fmt"
	agentModel "github.com/davyxu/cellmesh/svc/node/agent/model"
	"github.com/davyxu/cellmesh/svc/node/agent/routerule"
	"github.com/davyxu/protoplus/model"
	"github.com/davyxu/protoplus/util"
	"github.com/davyxu/ulog"
	"os"
)

// 从Proto文件中获取路由信息
func GenRouteTable(dset *model.DescriptorSet) (ret *agentModel.RouteTable) {

	ret = new(agentModel.RouteTable)

	for _, d := range dset.Structs() {

		msgDir := ParseMessage(d)
		msgID := msgidutil.StructMsgID(d)

		if msgDir.Valid() {
			var relay string
			if msgDir.Mid != "" {
				relay = fmt.Sprintf("[%s]", msgDir.Mid)
			}

			ulog.Debugf("  %s(%d)  %s -> %s   %s", d.Name, msgID, msgDir.From, msgDir.To, relay)

			ret.Rule = append(ret.Rule, &agentModel.RouteRule{
				MsgName: d.Name,
				SvcName: msgDir.To,
				MsgID:   msgID,
			})
		}

	}

	return
}

var (
	flagConfigKey = flag.String("configkey", agentModel.ConfigKey, "discovery kv config path")
	flagPackage   = flag.String("package", "proto", "package name in source files")
)

func main() {

	flag.Parse()

	ulog.SetLevel(ulog.DebugLevel)

	sd := redsd.NewRedisDiscovery()
	sd.Start("localhost:6379")

	dset := new(model.DescriptorSet)
	dset.PackageName = *flagPackage

	var routeTable *agentModel.RouteTable

	err := util.ParseFileList(dset)

	if err != nil {
		goto OnError
	}

	routeTable = GenRouteTable(dset)

	err = routerule.Upload(sd, routeTable, *flagConfigKey)

	if err != nil {
		goto OnError
	}

	return

OnError:
	fmt.Println(err)
	os.Exit(1)
}
