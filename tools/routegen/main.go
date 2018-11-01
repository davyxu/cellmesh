package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/davyxu/cellmesh/demo/table"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/discovery/consul"
	"github.com/davyxu/golog"
	"github.com/davyxu/protoplus/model"
	"github.com/davyxu/protoplus/util"
	"os"
)

// 从Proto文件中获取路由信息
func GenRouteTable(dset *model.DescriptorSet) (ret *table.RouteTable) {

	ret = new(table.RouteTable)

	for _, d := range dset.Structs() {

		if d.TagValueString("RouteRule") != "" && d.TagValueString("Service") != "" {

			ret.Rule = append(ret.Rule, &table.RouteRule{
				MsgName: d.Name,
				SvcName: d.TagValueString("Service"),
				Mode:    d.TagValueString("RouteRule"),
			})
		}
	}

	return
}

// 上传路由表到consul KV
func UploadRouteTable(tab *table.RouteTable) error {

	data, err := json.MarshalIndent(tab, "", "\t")

	if err != nil {
		return err
	}

	fmt.Printf("Write '%s'", *flagConfigPath)
	return discovery.Default.SetValue(*flagConfigPath, string(data))
}

var (
	flagConfigPath = flag.String("configpath", "config_demo/route_rule", "consul kv config path")
)

func main() {

	flag.Parse()

	discovery.Default = consulsd.NewDiscovery(nil)

	golog.SetLevelByString("consul", "info")

	dset := new(model.DescriptorSet)

	var routeTable *table.RouteTable

	err := util.ParseFileList(dset)

	if err != nil {
		goto OnError
	}

	routeTable = GenRouteTable(dset)

	err = UploadRouteTable(routeTable)

	if err != nil {
		goto OnError
	}

	return

OnError:
	fmt.Println(err)
	os.Exit(1)
}
