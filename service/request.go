package service

import (
	"errors"
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"reflect"
)

func selectStrategy(descList []*discovery.ServiceDesc) *discovery.ServiceDesc {

	if len(descList) == 0 {
		return nil
	}

	return descList[0]
}

func Request(serviceName string, req interface{}, ackType reflect.Type, callback func(interface{})) error {

	descList, err := discovery.Default.Query(serviceName)
	if err != nil {
		return err
	}

	desc := selectStrategy(descList)

	if desc == nil {
		return errors.New("target not reachable")
	}

	log.Debugf("Select service, %s", desc.String())

	addr := fmt.Sprintf("%s:%d", desc.Address, desc.Port)

	if rawConn, ok := connByAddr.Load(addr); ok {
		conn := rawConn.(Requestor)

		return conn.Request(req, ackType, callback)
	}

	return errors.New("connection not ready")
}
