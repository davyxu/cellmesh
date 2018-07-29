package svc

import (
	"errors"
	"fmt"
	"github.com/davyxu/cellmicro/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"reflect"
	"sync"
)

func Request(req, ack interface{}) error {

	svc, err := discovery.Default.Query("cellmicro.greating")
	if err != nil {
		return err
	}

	if len(svc) == 0 {
		return errors.New("target not reachable")
	}

	ackType := reflect.TypeOf(ack)

	if ackType.Kind() == reflect.Ptr {
		ackType = ackType.Elem()
	} else {
		return errors.New("invalid ack type, require ptr")
	}

	var waitMsg sync.WaitGroup
	waitMsg.Add(1)

	p := peer.NewGenericPeer("tcp.Connector", "node", fmt.Sprintf("%s:%d", svc[0].Address, svc[0].Port), nil)

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {

		incomingMsgType := reflect.TypeOf(ev.Message())
		if incomingMsgType.Elem() == ackType {
			ack = ev.Message()
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
