package routerule

import (
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/discovery"
)

// 用Consul KV下载路由规则
func Download() error {

	log.Debugf("Download route rule from discovery...")

	var tab model.RouteTable

	err := discovery.Default.GetValue(model.ConfigPath, &tab)
	if err != nil {
		return err
	}

	model.ClearRule()

	for _, r := range tab.Rule {
		model.AddRouteRule(r)
	}

	log.Debugf("Total %d rules added", len(tab.Rule))

	return nil
}
