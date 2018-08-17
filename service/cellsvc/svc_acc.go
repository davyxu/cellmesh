package cellsvc

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx/model"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

type accService struct {
	svcName string
	dis     service.DispatcherFunc

	listener cellnet.GenericPeer
}

func (self *accService) SetDispatcher(dis service.DispatcherFunc) {

	self.dis = dis
}

func (self *accService) Start() {

	self.listener = peer.NewGenericPeer("tcp.Acceptor", self.svcName, ":0", nil)

	proc.BindProcessorHandler(self.listener, "tcp.ltv", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *proto.ServiceIdentifyACK:

			if pre := service.GetConn(msg.SvcID); pre == nil {

				service.AddConn(ev.Session(), &discovery.ServiceDesc{
					ID:   msg.SvcID,
					Name: msg.SvcName,
				})
			}

		case *cellnet.SessionClosed:
			service.RemoveConn(ev.Session())
		}

		if self.dis != nil {

			self.dis(&svcEvent{
				Event: ev,
			})
		}
	})

	self.listener.Start()

	host := util.GetLocalIP()

	sd := &discovery.ServiceDesc{
		Host: host,
		Port: self.listener.(cellnet.TCPAcceptor).Port(),
		ID:   fxmodel.GetSvcID(self.svcName),
		Name: self.svcName,
	}

	log.SetColor("green").Debugf("service '%s' listen at %s:%d", sd.ID, host, sd.Port)

	discovery.Default.Register(sd)
}

func (self *accService) Stop() {
	discovery.Default.Deregister(fxmodel.GetSvcID(self.svcName))
}

func NewService(svcName string) service.Service {

	return &accService{
		svcName: svcName,
	}
}
