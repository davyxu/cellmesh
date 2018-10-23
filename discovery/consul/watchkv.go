package consulsd

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
)

func (self *consulDiscovery) startWatchKV() {

	plan, err := watch.Parse(map[string]interface{}{
		"type":   "keyprefix",
		"prefix": "/",
	})
	if err != nil {
		log.Errorln("startWatchKV:", err)
		return
	}

	plan.Handler = self.onKVListChanged

	go plan.Run(self.config.Address)

}

type KVMeta struct {
	Value []byte
	Plan  *watch.Plan
}

func (self *consulDiscovery) onKVListChanged(u uint64, data interface{}) {

	kvNames, ok := data.(api.KVPairs)
	if !ok {
		return
	}

	for _, kv := range kvNames {

		// 已经在cache里的,肯定添加过watch了
		if _, ok := self.kvCache.Load(kv.Key); ok {
			continue
		}

		plan, err := watch.Parse(map[string]interface{}{
			"type": "key",
			"key":  kv.Key,
		})

		if err == nil {
			plan.Handler = self.onKVChanged
			go plan.Run(self.config.Address)

			//log.Debugf("add kv : '%s'", kv.Key)

			self.kvCache.Store(kv.Key, &KVMeta{
				Value: kv.Value,
				Plan:  plan,
			})
		}
	}

	var foundKey string

	for {

		self.kvCache.Range(func(key, value interface{}) bool {

			kvKey := key.(string)
			meta := value.(*KVMeta)

			if !existsInPairs(kvNames, kvKey) {
				meta.Plan.Stop()
				foundKey = kvKey
				return false
			}

			return true
		})

		if foundKey == "" {
			break
		}

		//log.Debugf("remove kv : '%s'", foundKey)

		self.kvCache.Delete(foundKey)

		foundKey = ""
	}

}

func existsInPairs(kvp api.KVPairs, key string) bool {

	for _, kv := range kvp {
		if kv.Key == key {
			return true
		}
	}

	return false
}

func (self *consulDiscovery) onKVChanged(u uint64, data interface{}) {
	kv, ok := data.(*api.KVPair)
	if !ok {
		return
	}

	//log.Debugf("modify kv : '%s'", kv.Key)

	if raw, ok := self.kvCache.Load(kv.Key); ok {
		raw.(*KVMeta).Value = kv.Value
	}

}
