package service

import (
	"errors"
	"fmt"
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

// 保持长连接
func KeepConnection(reqSpawner func(addr string) Requestor, addr string, onReady chan Requestor) {
	requestor := reqSpawner(addr)

	requestor.Start()

	if requestor.IsReady() {

		connByAddr.Store(addr, requestor)
		log.SetColor("green").Debugln("add connection: ", addr)

		if onReady != nil {
			onReady <- requestor
		}

		// 连接断开
		requestor.WaitStop()
		connByAddr.Delete(addr)

		log.SetColor("yellow").Debugln("connection removed: ", addr)
	} else {

		requestor.Stop()
		time.Sleep(time.Second * 3)
	}
}

// 建立短连接
func CreateConnection(serviceName string, reqSpawner func(addr string) Requestor) (Requestor, error) {

	addr, err := QueryServiceAddress(serviceName)
	if err != nil {
		return nil, err
	}

	requestor := reqSpawner(addr)

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
		addr, err := QueryServiceAddress(serviceName)

		if err == nil {
			KeepConnection(reqSpawner, addr, onReady)
		}

		<-notify
	}

}
