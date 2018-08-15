package cellsvc

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"sync"
	"time"
)

type conService struct {
	svcName string
	dis     *service.Dispatcher

	descMap sync.Map

	mainSD *discovery.ServiceDesc
}

func (self *conService) SetDispatcher(dis *service.Dispatcher) {

	self.dis = dis
}

type connector interface {
	cellnet.TCPConnector
	IsReady() bool
}

func (self *conService) connFlow(desc *discovery.ServiceDesc) {

	if _, ok := self.descMap.Load(desc.Address()); ok {
		return
	}

	self.descMap.Store(desc.Address(), desc)

	var stop sync.WaitGroup

	p := peer.NewGenericPeer("tcp.SyncConnector", self.svcName, desc.Address(), nil)
	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionConnected:
			ev.Session().Send(proto.ServiceIdentifyACK{
				SvcName: self.mainSD.Name,
				SvcID:   self.mainSD.ID,
				Host:    self.mainSD.Host,
				Port:    int32(self.mainSD.Port),
			})

		case *cellnet.SessionClosed:
			stop.Done()
		}

		if self.dis != nil {
			self.dis.Invoke(ev)
		}
	})

	stop.Add(1)

	p.Start()

	conn := p.(connector)

	if conn.IsReady() {

		if desc != nil {

			service.AddConn(conn.Session(), desc)
		}

		// 连接断开
		stop.Wait()

		if desc != nil {
			service.RemoveConn(conn.Session())
		}

	} else {

		p.Stop()
		time.Sleep(time.Second * 3)
	}

	self.descMap.Delete(desc.Address())
}

func (self *conService) loop() {
	notify := discovery.Default.RegisterAddNotify()
	for {

		descList, err := discovery.Default.Query(self.svcName)
		if err == nil && len(descList) > 0 {

			// 保持服务发现中的所有连接
			for _, desc := range descList {

				go self.connFlow(desc)
			}

		}

		<-notify
	}
}

func (self *conService) Start() {

	go self.loop()
}

func (self *conService) Stop() {

}

func NewConnService(svcName string, mainSD *discovery.ServiceDesc) service.Service {

	return &conService{
		svcName: svcName,
		mainSD:  mainSD,
	}
}
