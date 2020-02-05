package routerule

import (
	"encoding/json"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/svc/agent/model"
	"github.com/davyxu/ulog"
)

// 用Consul KV下载路由规则
func Download() error {

	ulog.Infof("Download route rule from discovery...")

	var tab model.RouteTable

	err := discovery.Global.GetValue(model.ConfigKey, &tab)
	if err != nil {
		return err
	}

	model.ClearRule()

	for _, r := range tab.Rule {
		model.AddRouteRule(r)
	}

	ulog.Infof("Total %d route rules", len(tab.Rule))

	return nil
}

// 上传路由表到consul KV
func Upload(tab *model.RouteTable, configKey string) error {

	data, err := json.MarshalIndent(tab, "", "\t")

	if err != nil {
		return err
	}

	ulog.Infof("Write '%s', count: %d", configKey, len(tab.Rule))
	return discovery.Global.SetValue(configKey, string(data))
}
