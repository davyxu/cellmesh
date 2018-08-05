package service

import (
	"errors"
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"reflect"
	"sync"
)

type Requestor interface {
	Request(req interface{}, ackType reflect.Type, callback func(interface{})) error
}

var (
	connByAddr sync.Map

	NewRequestor func(addr string, readyChan chan Requestor) Requestor
)

func RemoveConnection(addr string) {
	connByAddr.Delete(addr)
}

func PrepareConnection(serviceName string) error {

	descList, err := discovery.Default.Query(serviceName)
	if err != nil {
		return err
	}

	desc := selectStrategy(descList)

	if desc == nil {
		return errors.New("target not reachable")
	}

	addr := fmt.Sprintf("%s:%d", desc.Address, desc.Port)

	readyConn := make(chan Requestor)

	if NewRequestor == nil {
		panic("requestor package not import")
	}

	NewRequestor(addr, readyConn)

	connByAddr.Store(addr, <-readyConn)

	return nil
}
