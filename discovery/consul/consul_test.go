package consulsd

import (
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"math/rand"
	"testing"
	"time"
)

func regService(sd discovery.Discovery) {

	svcid := fmt.Sprintf("svc%d", rand.Int31n(10))

	svcname := fmt.Sprintf("type%d", rand.Int31n(3))

	sd.Register(&discovery.ServiceDesc{
		Name: svcname,
		ID:   svcid,
	})
}

func unregService(sd discovery.Discovery) {

	svcid := fmt.Sprintf("svc%d", rand.Int31n(10))

	sd.Deregister(svcid)
}

func TestConsul(t *testing.T) {

	sd := NewDiscovery(nil)

	for {
		go regService(sd)
		time.Sleep(time.Millisecond * 10)
		go unregService(sd)
	}

	for {
		svcname := fmt.Sprintf("type%d", rand.Int31n(3))

		for _, desc := range sd.Query(svcname) {
			t.Log(desc.String())
		}

		if len(sd.Query(svcname)) == 0 {
			t.Errorf("query nil")
		}
	}
}
