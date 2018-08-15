package cellsvc

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	meshutil "github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

type accService struct {
	svcName    string
	listenPort int
	dis        *service.Dispatcher

	sd *discovery.ServiceDesc
}

func (self *accService) ID() string {
	return fmt.Sprintf("%s-%d", self.svcName, self.listenPort)
}

func (self *accService) SetDispatcher(dis *service.Dispatcher) {

	self.dis = dis
}

func (self *accService) Start() {

	p := peer.NewGenericPeer("tcp.Acceptor", self.svcName, ":0", nil)

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *proto.ServiceIdentifyACK:

			if pre := service.GetConn(msg.SvcID); pre == nil {

				service.AddConn(ev.Session(), &discovery.ServiceDesc{
					ID:   msg.SvcID,
					Name: msg.SvcName,
					Host: msg.Host,
					Port: int(msg.Port),
				})
			}

		case *cellnet.SessionClosed:
			service.RemoveConn(ev.Session())
		}

		if self.dis != nil {
			self.dis.Invoke(ev)
		}
	})

	p.Start()

	self.listenPort = p.(cellnet.TCPAcceptor).ListenPort()

	host := meshutil.GetLocalIP()

	self.sd = &discovery.ServiceDesc{
		Host: host,
		Port: self.listenPort,
		ID:   self.ID(),
		Name: self.svcName,
	}

	log.SetColor("green").Debugf("service '%s' listen at %s:%d", self.svcName, host, self.listenPort)

	discovery.Default.Register(self.sd)
}

func (self *accService) GetSD() *discovery.ServiceDesc {
	return self.sd
}

func (self *accService) Stop() {
	discovery.Default.Deregister(self.ID())
}

func NewService(svcName string) service.Service {

	return &accService{
		svcName: svcName,
	}
}
