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

func GetSession(addr string) cellnet.Session {
	if rawConn, ok := connByAddr.Load(addr); ok {
		conn := rawConn.(Requestor)

		return conn.Session()
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
func KeepConnection(requestor Requestor, svcid string, onReady chan Requestor) {

	requestor.Start()

	if requestor.IsReady() {

		if svcid != "" {
			connByAddr.Store(svcid, requestor)
			log.SetColor("green").Debugln("add connection: ", svcid)
		}

		if onReady != nil {
			onReady <- requestor
		}

		// 连接断开
		requestor.WaitStop()
		if svcid != "" {
			connByAddr.Delete(svcid)

			log.SetColor("yellow").Debugln("connection removed: ", svcid)
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
func PrepareConnection(serviceName string, reqSpawner func(addr string) Requestor, onReady chan Requestor) {

	notify := discovery.Default.RegisterAddNotify()
	for {
		desc, err := QueryServiceAddress(serviceName)

		if err == nil {
			KeepConnection(reqSpawner(desc.Address()), desc.ID, onReady)
		}

		<-notify
	}

}
