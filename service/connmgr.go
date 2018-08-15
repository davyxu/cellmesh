package service

import (
	"errors"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"reflect"
	"sync"
	"time"
)

type Requestor interface {
	Request(req interface{}, ackType reflect.Type, callback func(interface{})) error

	Session() cellnet.Session

	Start()

	IsReady() bool

	Stop()

	WaitStop()
}

var (
	connByAddr sync.Map
)

func GetSession(svcid string) cellnet.Session {

	if raw, ok := connByAddr.Load(svcid); ok {
		return raw.(cellnet.Session)
	}

	return nil
}

func QueryServiceAddress(serviceName string) (*discovery.ServiceDesc, error) {
	descList, err := discovery.Default.Query(serviceName)
	if err != nil {
		return nil, err
	}

	desc := selectStrategy(descList)

	if desc == nil {
		return nil, errors.New("target not reachable")
	}

	return desc, nil
}

// 保持长连接
func KeepConnection(requestor Requestor, desc *discovery.ServiceDesc, onReady func(*discovery.ServiceDesc, Requestor)) {

	requestor.Start()

	if requestor.IsReady() {

		if desc != nil {
			connByAddr.Store(desc.ID, requestor.Session())
			log.SetColor("green").Debugln("add connection: ", desc.ID)
		}

		if onReady != nil {
			onReady(desc, requestor)
		}

		// 连接断开
		requestor.WaitStop()
		if desc != nil {
			connByAddr.Delete(desc.ID)

			log.SetColor("yellow").Debugln("connection removed: ", desc.ID)
		}

	} else {

		requestor.Stop()
		time.Sleep(time.Second * 3)
	}
}

// 建立短连接
func CreateConnection(serviceName string, reqSpawner func(addr string) Requestor) (Requestor, error) {

	desc, err := QueryServiceAddress(serviceName)
	if err != nil {
		return nil, err
	}

	requestor := reqSpawner(desc.Address())

	requestor.Start()

	if requestor.IsReady() {
		return requestor, err
	}

	return nil, errors.New("fail create connection")
}

// 异步建立与服务连接
func PrepareConnection(serviceName string, reqSpawner func(addr string) Requestor, onReady func(*discovery.ServiceDesc, Requestor)) {

	notify := discovery.Default.RegisterAddNotify()
	for {
		desc, err := QueryServiceAddress(serviceName)

		if err == nil {
			KeepConnection(reqSpawner(desc.Address()), desc, onReady)
		}

		<-notify
	}

}
