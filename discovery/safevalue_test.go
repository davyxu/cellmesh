package discovery

import (
	"github.com/davyxu/cellmesh/discovery/consul"
	"github.com/davyxu/cellmesh/service"
	"reflect"
	"testing"
)

func TestSafeGetValue(t *testing.T) {

	var origin []byte
	for i := 0; i < 12; i++ {
		//origin = append(origin, byte(rand.Int31n(127)))
		origin = append(origin, byte(i))
	}

	sdConfig := consulsd.DefaultConfig()
	sdConfig.Address = service.GetDiscoveryAddr()
	Default = consulsd.NewDiscovery(sdConfig)

	err := SafeSetValue(Default, "config/test", origin, true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var outData []byte
	err = SafeGetValue(Default, "config/test", &outData, true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(origin, outData) {
		t.FailNow()
	}
}
