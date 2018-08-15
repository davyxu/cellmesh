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
	connBySvcID        = map[string]cellnet.Session{}
	connBySvcNameGuard sync.RWMutex
)

func AddConn(ses cellnet.Session, desc *discovery.ServiceDesc) {

	connBySvcNameGuard.Lock()
	ses.(cellnet.ContextSet).SetContext("desc", desc)
	connBySvcID[desc.ID] = ses
	connBySvcNameGuard.Unlock()

	log.SetColor("green").Debugln("add connection: ", desc.ID)
}

func GetSDBySession(ses cellnet.Session) (ret *discovery.ServiceDesc) {

	connBySvcNameGuard.RLock()
	defer connBySvcNameGuard.RUnlock()

	for _, libses := range connBySvcID {
		if libses == ses {
			ses.(cellnet.ContextSet).GetContext("desc", &ret)
			break
		}
	}

	return
}

func GetConn(svcid string) cellnet.Session {
	connBySvcNameGuard.RLock()
	defer connBySvcNameGuard.RUnlock()

	if ses, ok := connBySvcID[svcid]; ok {

		return ses
	}

	return nil
}

func RemoveConn(ses cellnet.Session) {
	var desc *discovery.ServiceDesc
	if ses.(cellnet.ContextSet).GetContext("desc", &desc) {
		connBySvcNameGuard.Lock()
		delete(connBySvcID, desc.ID)
		connBySvcNameGuard.Unlock()

		log.SetColor("yellow").Debugln("connection removed: ", desc.ID)
	}
}

func VisitConn(callback func(ses cellnet.Session, desc *discovery.ServiceDesc)) {
	connBySvcNameGuard.RLock()

	for _, ses := range connBySvcID {

		var desc *discovery.ServiceDesc
		if ses.(cellnet.ContextSet).GetContext("desc", &desc) {
			callback(ses, desc)
		}
	}

	connBySvcNameGuard.RUnlock()
}

func QueryServiceAddress(serviceName string) (*discovery.ServiceDesc, error) {
	descList, err := discovery.Default.Query(serviceName)
	if err != nil {
		return nil, err
	}

	desc := selectStrategy(descList)

	if desc == nil {
		return nil, errors.New("target not reachable:" + serviceName)
	}

	return desc, nil
}

// 建立短连接
func CreateConnection(serviceName string, reqSpawner func(addr string) Requestor) (Requestor, error) {

	notify := discovery.Default.RegisterNotify("add")
	for {

		desc, err := QueryServiceAddress(serviceName)

		if err == nil {

			requestor := reqSpawner(desc.Address())

			requestor.Start()

			if requestor.IsReady() {
				return requestor, err
			}

			requestor.Stop()
		}

		<-notify
	}

	discovery.Default.DeregisterNotify("add", notify)

	return nil, nil
}

// 保持长连接
func KeepConnection(requestor Requestor, svcid string, onReady chan Requestor) {

	requestor.Start()

	if requestor.IsReady() {

		if svcid != "" {
			log.SetColor("green").Debugln("add connection: ", svcid)
		}

		if onReady != nil {
			onReady <- requestor
		}

		// 连接断开
		requestor.WaitStop()
		if svcid != "" {
			log.SetColor("yellow").Debugln("connection removed: ", svcid)
		}

	} else {

		requestor.Stop()
		time.Sleep(time.Second * 3)
	}
}
