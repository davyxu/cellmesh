package zonecfg

import (
	"github.com/davyxu/ulog"
	"sync"
)

var (
	Tab        = NewTable()
	currZone   string
	currNode   string
	valueByKey sync.Map
)

func matchNode(def *ZoneConfig) bool {
	if len(def.NodeList) > 0 {

		for _, node := range def.NodeList {
			if node == "*" {
				return true
			}

			if node == currNode {
				return true
			}
		}

		return false

	} else {
		ulog.Warnf("Empty node list, ignore, %+v", *def)
		return false
	}
}

func init() {

	Tab.RegisterPostEntry(func(tab *Table) error {

		for index, def := range tab.ZoneConfig {

			// 只加载当前区配置
			if def.Zone == currZone && matchNode(def) {

				if _, ok := valueByKey.Load(def.Key); ok {
					ulog.Warnf("Duplicate ZoneConfig, %+v", *def)
					continue
				}

				def.Line = int32(index)
				valueByKey.Store(def.Key, def)
			}
		}

		return nil
	})
}
