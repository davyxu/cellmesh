package consulsd

import "github.com/hashicorp/consul/api"

func (self *consulDiscovery) SetValue(key string, value []byte) error {

	_, err := self.client.KV().Put(&api.KVPair{
		Key:   key,
		Value: value,
	}, nil)

	return err
}

func (self *consulDiscovery) GetValue(key string) ([]byte, bool, error) {

	kvPair, _, err := self.client.KV().Get(key, nil)

	if err != nil {
		return nil, false, err
	}

	// 值不存在
	if kvPair == nil {
		return nil, false, nil
	}

	return kvPair.Value, true, nil
}
