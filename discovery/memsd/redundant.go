package main

import (
	"github.com/davyxu/cellmesh/discovery/memsd/model"
	"time"
)

// 移除token丢失的values
func startCheckRedundantValue() {

	ticker := time.NewTicker(time.Minute)

	for {

		<-ticker.C

		// 与收发在一个队列中，保证无锁
		model.Queue.Post(func() {

			var svcToDelete []*model.ValueMeta

			model.VisitValue(func(meta *model.ValueMeta) bool {

				if meta.Token != "" && !model.TokenExists(meta.Token) && meta.ValueAsServiceDesc().GetMeta("@Persist") == "" {
					svcToDelete = append(svcToDelete, meta)
				}

				return true
			})

			for _, meta := range svcToDelete {
				deleteNotify(meta.Key, "check redundant")
			}
		})
	}

}
