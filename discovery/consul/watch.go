package consulsd

import (
	"github.com/hashicorp/consul/watch"
)

func (self *consulDiscovery) startWatch() {

	wp, err := watch.Parse(map[string]interface{}{"type": "service", "service": "cellmicro.greating2"})
	if err != nil {
		log.Errorln(err)
		return
	}

	wp.Handler = self.onNameListChanged

	go wp.Run(self.config.Address)

}

func (self *consulDiscovery) onNameListChanged(u uint64, data interface{}) {
	svcNames, ok := data.(map[string][]string)
	if !ok {
		return
	}

	for svcName := range svcNames {

		// 已经对这种名称的服务创建了watcher的跳过
		if _, ok := self.nameWatcher[svcName]; ok {
			continue
		}

		// 发现新的服务
		wp, err := watch.Parse(map[string]interface{}{
			"type":    "service",
			"service": svcName,
		})

		if err == nil {
			wp.Handler = self.onServiceChanged
			go wp.Run(self.config.Address)
			self.nameWatcher[svcName] = wp

		}
	}
}

func (self *consulDiscovery) onServiceChanged(u uint64, data interface{}) {
	//svcDetails, ok := data.([]*api.ServiceEntry)
	//if !ok {
	//	return
	//}
	//
	//for _, detail := range svcDetails {
	//
	//}
}
