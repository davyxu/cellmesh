package endpoint

import (
	"errors"
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"reflect"
	"sync"
)

func Request(req interface{}, ackType reflect.Type, callback func(interface{})) error {

	svc, err := discovery.Default.Query("cellmicro.greating")
	if err != nil {
		return err
	}

	if len(svc) == 0 {
		return errors.New("target not reachable")
	}

	var waitMsg sync.WaitGroup
	waitMsg.Add(1)

	p := peer.NewGenericPeer("tcp.Connector", "node", fmt.Sprintf("%s:%d", svc[0].Address, svc[0].Port), nil)

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {

		incomingMsgType := reflect.TypeOf(ev.Message())
		if incomingMsgType.Elem() == ackType {
			callback(ev.Message())
			waitMsg.Done()
		}

		switch ev.Message().(type) {
		case *cellnet.SessionConnected: // 已经连接上
			ev.Session().Send(req)
		}
	})

	p.Start()

	waitMsg.Wait()

	return nil
}
