package discovery

import (
	"testing"
	"github.com/davyxu/cellmesh/discovery/consul"
)

func Test_Discovery(t *testing.T) {
	sdConfig := consulsd.DefaultConfig()
	sdConfig.Address = "127.0.0.1:8500"
	Default = consulsd.NewDiscovery(sdConfig)


}
