package service

import (
	"errors"
	"github.com/davyxu/cellmesh/discovery"
	"reflect"
)

func selectStrategy(descList []*discovery.ServiceDesc) *discovery.ServiceDesc {

	if len(descList) == 0 {
		return nil
	}

	return descList[0]
}

var (
	ErrInvalidTarget = errors.New("target provider should be 'servicename' or 'Requestor'")
)

func Request(targetProvider interface{}, req interface{}, ackType reflect.Type, callback func(interface{})) error {

	var requestor Requestor
	switch tgt := targetProvider.(type) {
	case string:
		addr, err := QueryServiceAddress(tgt)
		if err != nil {
			return err
		}

		if rawConn, ok := connByAddr.Load(addr); ok {
			requestor = rawConn.(Requestor)
		}

	case Requestor:
		requestor = tgt
	default:
		panic(ErrInvalidTarget)
	}

	return requestor.Request(req, ackType, callback)
}
