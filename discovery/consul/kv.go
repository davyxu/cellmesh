package consulsd

import "github.com/hashicorp/consul/api"

func (self *consulDiscovery) SetValue(key string, value []byte) error {

	_, err := self.client.KV().Put(&api.KVPair{
		Key:   key,
		Value: value,
	}, nil)

	return err
}

func (self *consulDiscovery) GetValue(key string) ([]byte, error) {

	kvPair, _, err := self.client.KV().Get(key, nil)

	if err != nil || kvPair == nil {
		return nil, err
	}

	return kvPair.Value, nil
}
