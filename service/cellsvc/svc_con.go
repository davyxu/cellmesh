package cellsvc

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx/model"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"sync"
	"time"
)

type connector interface {
	cellnet.TCPConnector
	IsReady() bool
}

type conService struct {
	evDispatcher

	svcName string

	connectorBySvcID sync.Map // map[svcid] connector
}

func (self *conService) connFlow(p cellnet.GenericPeer, sd *discovery.ServiceDesc) {

	var stop sync.WaitGroup

	proc.BindProcessorHandler(p, self.procName, func(ev cellnet.Event) {

		switch ev.Message().(type) {
		case *cellnet.SessionConnected:
			ev.Session().Send(proto.ServiceIdentifyACK{
				SvcName: self.svcName,
				SvcID:   fxmodel.GetSvcID(self.svcName),
			})

		case *cellnet.SessionClosed:
			stop.Done()
		}

		self.Invoke(&svcEvent{
			Event: ev,
		})
	})

	stop.Add(1)

	p.Start()

	conn := p.(connector)

	if conn.IsReady() {

		if sd != nil {

			service.AddRemoteService(conn.Session(), sd)
		}

		// 连接断开
		stop.Wait()

		if sd != nil {
			service.RemoveRemoteService(conn.Session())
		}

	} else {

		p.Stop()
		time.Sleep(time.Second * 3)
	}

	self.connectorBySvcID.Delete(sd.ID)
}

func (self *conService) loop() {
	notify := discovery.Default.RegisterNotify("add")
	for {

		log.Debugf("Query svc from discovery, name: '%s'...", self.svcName)
		descList, err := discovery.Default.Query(self.svcName)
		if err == nil && len(descList) > 0 {

			// 保持服务发现中的所有连接
			for _, sd := range descList {

				// 新连接马上连接，老连接保留
				if _, ok := self.connectorBySvcID.Load(sd.ID); !ok {

					p := peer.NewGenericPeer("tcp.SyncConnector", self.svcName, sd.Address(), nil)
					self.connectorBySvcID.Store(sd.ID, p)

					go self.connFlow(p, sd)
				}
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

func NewConnector(svcName string) service.Service {

	return &conService{
		svcName: svcName,
	}
}
