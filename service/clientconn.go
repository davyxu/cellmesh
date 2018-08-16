package service

import (
	"errors"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/tcp"
	"reflect"
	"sync"
	"time"
)

// 建立短连接
func CreateConnection(serviceName string) (cellnet.Session, error) {

	notify := discovery.Default.RegisterNotify("add")
	for {

		desc, err := QueryServiceAddress(serviceName)

		if err == nil {

			p := peer.NewGenericPeer("tcp.SyncConnector", serviceName, desc.Address(), nil)
			proc.BindProcessorHandler(p, "cellmesh.tcp", nil)

			p.Start()

			conn := p.(connector)

			if conn.IsReady() {
				return conn.Session(), err
			}

			p.Stop()
		}

		<-notify
	}

	discovery.Default.DeregisterNotify("add", notify)

	return nil, nil
}

type connector interface {
	cellnet.TCPConnector
	IsReady() bool
}

// 保持长连接
func KeepConnection(svcid, addr string, onReady chan cellnet.Session) {

	var stop sync.WaitGroup

	p := peer.NewGenericPeer("tcp.SyncConnector", svcid, addr, nil)
	proc.BindProcessorHandler(p, "cellmesh.tcp", func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionClosed:
			stop.Done()
		}
	})

	stop.Add(1)

	p.Start()

	conn := p.(connector)

	if conn.IsReady() {

		if onReady != nil {
			onReady <- conn.Session()
		}

		// 连接断开
		stop.Wait()

	} else {

		p.Stop()
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

var (
	RPCPairQueryFunc func(req interface{}) reflect.Type
)

func RemoteCall(target, req interface{}, callback func(interface{})) error {

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

	if RPCPairQueryFunc == nil {
		panic("Require 'RPCPairQueryFunc' from protogen code")
	}

	ackType := RPCPairQueryFunc(req)

	callByType.Store(ackType, feedBack)

	defer callByType.Delete(ackType)

	ses.Send(req)

	select {
	case ack := <-feedBack:
		callback(ack)
		return nil
	case <-time.After(time.Second):

		log.Errorln("RemoteCall: RPC time out")

		return errors.New("RPC Time out")
	}

	return nil
}

func init() {
	transmitter := new(tcp.TCPMessageTransmitter)
	typeRPCHooker := new(TypeRPCHooker)
	msgLogger := new(tcp.MsgHooker)

	proc.RegisterProcessor("cellmesh.tcp", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(transmitter)
		bundle.SetHooker(proc.NewMultiHooker(msgLogger, typeRPCHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
