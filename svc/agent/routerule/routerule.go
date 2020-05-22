package routerule

import (
	"encoding/json"
	"github.com/davyxu/cellmesh/link"
	"github.com/davyxu/cellmesh/redsd"
	"github.com/davyxu/cellmesh/svc/agent/model"
	"github.com/davyxu/ulog"
)

// 用Consul KV下载路由规则
func Download() {

	ulog.Infof("Download route rule from discovery...")

	var tab model.RouteTable

	err := link.SD.GetValue(model.ConfigKey, &tab)
	if err != nil {
		ulog.Errorf("route table get failed, %s", err)
		return
	}

	model.ClearRule()

	for _, r := range tab.Rule {
		model.AddRouteRule(r)
	}

	ulog.Infof("Total %d route rules", len(tab.Rule))
}

// 上传路由表到consul KV
func Upload(sd *redsd.RedisDiscovery, tab *model.RouteTable, configKey string) error {

	data, err := json.Marshal(tab)

	if err != nil {
		return err
	}

	ulog.Infof("Write '%s', count: %d", configKey, len(tab.Rule))
	return sd.SetValue(configKey, string(data))
}
