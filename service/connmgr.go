package service

import (
	"errors"
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"reflect"
	"sync"
)

type Requestor interface {
	Request(req interface{}, ackType reflect.Type, callback func(interface{})) error

	Session() cellnet.Session
}

var (
	connByAddr sync.Map

	NewRequestor func(addr string, readyChan chan Requestor) Requestor
)

func GetSession(addr string) cellnet.Session {
	if rawConn, ok := connByAddr.Load(addr); ok {
		conn := rawConn.(Requestor)

		return conn.Session()
	}

	return nil
}

func RemoveConnection(addr string) {
	connByAddr.Delete(addr)
}

func QueryServiceAddress(serviceName string) (string, error) {
	descList, err := discovery.Default.Query(serviceName)
	if err != nil {
		return "", err
	}

	desc := selectStrategy(descList)

	if desc == nil {
		return "", errors.New("target not reachable")
	}

	return fmt.Sprintf("%s:%d", desc.Address, desc.Port), nil
}

func PrepareConnection(serviceName string) error {

	addr, err := QueryServiceAddress(serviceName)
	if err != nil {
		return err
	}

	readyConn := make(chan Requestor)

	if NewRequestor == nil {
		panic("requestor package not import")
	}

	NewRequestor(addr, readyConn)

	connByAddr.Store(addr, <-readyConn)

	return nil
}
