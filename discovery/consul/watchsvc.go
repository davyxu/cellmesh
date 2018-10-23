package consulsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
)

func (self *consulDiscovery) startWatchService() {

	plan, err := watch.Parse(map[string]interface{}{"type": "services"})
	if err != nil {
		log.Errorln("startWatchService:", err)
		return
	}

	plan.Handler = self.onSvcNameListChanged

	go plan.Run(self.config.Address)

}

func (self *consulDiscovery) onSvcNameListChanged(u uint64, data interface{}) {
	svcNames, ok := data.(map[string][]string)
	if !ok {
		return
	}

	for svcName := range svcNames {

		// 已经对这种名称的服务创建了watcher的跳过
		if _, ok := self.svcNameWatcher.Load(svcName); ok {
			continue
		}

		// 发现新的服务
		plan, err := watch.Parse(map[string]interface{}{
			"type":    "service",
			"service": svcName,
		})

		if err == nil {
			plan.Handler = self.onServiceChanged
			go plan.Run(self.config.Address)

			//log.Debugf("Watch service '%s' begin", svcName)

			self.svcNameWatcher.Store(svcName, plan)
		}
	}

	var foundSvc string

	for {

		self.svcNameWatcher.Range(func(key, value interface{}) bool {

			svcName := key.(string)
			plan := value.(*watch.Plan)

			if _, ok := svcNames[svcName]; !ok {

				plan.Stop()

				foundSvc = svcName

				// 删除后重新扫描，直到没有发现要删除的为止
				return false
			}

			return true
		})

		if foundSvc == "" {
			break
		}

		self.svcNameWatcher.Delete(foundSvc)

		if raw, ok := self.svcCache.Load(foundSvc); ok {
			for _, svc := range raw.([]*discovery.ServiceDesc) {
				self.OnCacheUpdated("remove", svc)
			}
		}

		// 删除这个名字的所有缓冲的服务
		self.svcCache.Delete(foundSvc)

		foundSvc = ""
	}

}

func (self *consulDiscovery) onServiceChanged(u uint64, data interface{}) {
	svcDetails, ok := data.([]*api.ServiceEntry)
	if !ok || len(svcDetails) == 0 {
		return
	}

	// 防止多次触发时，并发写入cache内列表时互相覆盖
	self.svcCacheGuard.Lock()
	defer self.svcCacheGuard.Unlock()

	svcName := svcDetails[0].Service.Service

	var newList []*discovery.ServiceDesc

	for _, detail := range svcDetails {
		if isServiceHealth(detail) {
			newList = append(newList, consulSvcToService(detail.Service))
		}

	}

	var oldList []*discovery.ServiceDesc
	if raw, ok := self.svcCache.Load(svcName); ok {
		oldList = raw.([]*discovery.ServiceDesc)
	}

	self.svcCache.Store(svcName, newList)

	for _, oldSvc := range oldList {

		// 在新的列表中没有找到老的id，表示服务被移除
		if !existsInServiceList(newList, oldSvc.ID) {
			self.OnCacheUpdated("remove", oldSvc)
		}
	}

	for _, newSvc := range newList {

		// 在老的列表中没有找到新的id，表示服务新增
		if !existsInServiceList(oldList, newSvc.ID) {
			self.OnCacheUpdated("add", newSvc)
		}
	}

}
