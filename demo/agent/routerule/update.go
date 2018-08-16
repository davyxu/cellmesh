package routerule

import (
	"encoding/json"
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/discovery"
)

func Download() error {

	log.Debugf("Download route rule from discovery...")

	data, err := discovery.Default.GetValue(model.ConfigPath)
	if err != nil {
		return err
	}

	var tab model.RouteTable

	err = json.Unmarshal(data, &tab)
	if err != nil {
		return err
	}

	model.ClearRule()

	for _, r := range tab.Rule {
		model.AddRouteRule(r)
	}

	return nil
}
