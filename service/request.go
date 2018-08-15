package service

import (
	"errors"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/rpc"
	"reflect"
	"time"
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

func Request(targetProvider interface{}, req interface{}, ackType reflect.Type, callback func(interface{})) (err error) {

	var requestor Requestor
	switch tgt := targetProvider.(type) {
	case Requestor:
		requestor = tgt
	case cellnet.Session:
		log.Debugln(1)
		ack, err := rpc.CallSync(tgt, req, time.Second*5)
		if err != nil {
			log.Debugln(2)
			return err
		}

		log.Debugln(3)

		callback(ack)
		log.Debugln(4)
		return nil

	default:
		panic(ErrInvalidTarget)
	}

	if requestor != nil {
		err = requestor.Request(req, ackType, callback)
	} else {
		err = errors.New("target not ready")
	}

	if err != nil {
		log.Errorf("Request failed, req: %s %s", cellnet.MessageToName(req), err.Error())
	}

	return err
}
