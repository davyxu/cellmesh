package service

import (
	"errors"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/gorillaws"
	"github.com/davyxu/cellnet/proc/tcp"
	"reflect"
	"sync"
	"time"
)

func selectStrategy(descList []*discovery.ServiceDesc) *discovery.ServiceDesc {

	if len(descList) == 0 {
		return nil
	}

	return descList[0]
}

func queryServiceAddress(serviceName string) (*discovery.ServiceDesc, error) {
	descList := discovery.Default.Query(serviceName)

	desc := selectStrategy(descList)

	if desc == nil {
		return nil, errors.New("target not reachable:" + serviceName)
	}

	return desc, nil
}

// 建立短连接
func CreateConnection(serviceName, netPeerType, netProcName string) (ret cellnet.Session) {

	notify := discovery.Default.RegisterNotify("add")

	done := make(chan struct{})
	go func() {
		for {

			desc, err := queryServiceAddress(serviceName)

			if err == nil {

				p := peer.NewGenericPeer(netPeerType, serviceName, desc.Address(), nil)
				proc.BindProcessorHandler(p, netProcName, nil)

				p.Start()

				conn := p.(connector)

				if conn.IsReady() {
					ret = conn.Session()

					break
				}

				p.Stop()
			}

			<-notify
		}

		discovery.Default.DeregisterNotify("add", notify)
		done <- struct{}{}
	}()

	<-done

	return
}

type connector interface {
	cellnet.TCPConnector
	cellnet.PeerReadyChecker
}

// 保持长连接
func KeepConnection(svcid, addr, netPeerType, netProc string, onReady func(cellnet.Session), onClose func()) {

	var stop sync.WaitGroup

	p := peer.NewGenericPeer(netPeerType, svcid, addr, nil)
	proc.BindProcessorHandler(p, netProc, func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionClosed:
			stop.Done()
		}
	})

	stop.Add(1)

	p.Start()

	conn := p.(connector)

	if conn.IsReady() {

		onReady(conn.Session())

		// 连接断开
		stop.Wait()
	}

	p.Stop()

	if onClose != nil {
		onClose()
	}

}

var (
	callByType sync.Map // map[reflect.Type]func(interface{})
)

type TypeRPCHooker struct {
}

func (TypeRPCHooker) OnInboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	outputEvent, _, err := ResolveInboundEvent(inputEvent)
	if err != nil {
		log.Errorln("rpc.ResolveInboundEvent", err)
		return
	}

	return outputEvent
}

func (TypeRPCHooker) OnOutboundEvent(inputEvent cellnet.Event) (outputEvent cellnet.Event) {

	return inputEvent
}

func ResolveInboundEvent(inputEvent cellnet.Event) (ouputEvent cellnet.Event, handled bool, err error) {
	incomingMsgType := reflect.TypeOf(inputEvent.Message()).Elem()

	if rawFeedback, ok := callByType.Load(incomingMsgType); ok {
		feedBack := rawFeedback.(chan interface{})
		feedBack <- inputEvent.Message()
		return inputEvent, true, nil
	}

	return inputEvent, false, nil
}

// callback =func(ack *YouMsgACK)
func RemoteCall(target, req interface{}, callback interface{}) error {

	funcType := reflect.TypeOf(callback)
	if funcType.Kind() != reflect.Func {
		panic("callback require 'func'")
	}

	var ses cellnet.Session
	switch tgt := target.(type) {
	case cellnet.Session:
		ses = tgt
	default:
		panic("rpc: Invalid peer type, require cellnet.Session")
	}

	if ses == nil {
		return errors.New("Empty session")
	}

	feedBack := make(chan interface{})

	// 获取回调第一个参数

	if funcType.NumIn() != 1 {
		panic("callback func param format like 'func(ack *YouMsgACK)'")
	}

	ackType := funcType.In(0)
	if ackType.Kind() != reflect.Ptr {
		panic("callback func param format like 'func(ack *YouMsgACK)'")
	}

	ackType = ackType.Elem()

	callByType.Store(ackType, feedBack)

	defer callByType.Delete(ackType)

	ses.Send(req)

	select {
	case ack := <-feedBack:

		vCall := reflect.ValueOf(callback)

		vCall.Call([]reflect.Value{reflect.ValueOf(ack)})
		return nil
	case <-time.After(time.Second):

		log.Errorln("RemoteCall: RPC time out")

		return errors.New("RPC Time out")
	}

	return nil
}

func init() {
	proc.RegisterProcessor("cellmesh.tcp", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(tcp.MsgHooker), new(TypeRPCHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})

	proc.RegisterProcessor("cellmesh.ws", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(new(gorillaws.WSMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(tcp.MsgHooker), new(TypeRPCHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
