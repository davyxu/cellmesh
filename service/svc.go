package service

import (
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"reflect"
	"sync"
)

type cellService struct {
	name       string
	listenPort int
	svcByName  sync.Map // map[reflect.Type]*endpoint.MethodInfo
}

func (self *cellService) AddCall(name string, svc *MethodInfo) {

	self.svcByName.Store(svc.RequestType, svc)
}

func (self *cellService) ID() string {
	return fmt.Sprintf("%s-%d", self.name, self.listenPort)
}

type ReplyEvent interface {
	Reply(msg interface{})
}

func (self *cellService) Start() error {

	p := peer.NewGenericPeer("tcp.Acceptor", "", ":0", nil)

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {

		switch evData := ev.(type) {
		case ReplyEvent:

			msgType := reflect.TypeOf(ev.Message()).Elem()

			if svcRaw, ok := self.svcByName.Load(msgType); ok {

				svc := svcRaw.(*MethodInfo)

				e := &Event{
					Request:  ev.Message(),
					Response: svc.NewResponse(),
				}

				svc.Handler(e)

				evData.Reply(e.Response)
			}
		}

	})

	p.Start()

	self.listenPort = p.(cellnet.TCPAcceptor).ListenPort()

	p.(cellnet.PeerProperty).SetName(fmt.Sprintf(":%d", self.listenPort))

	host := util.GetLocalIP()

	sd := &discovery.ServiceDesc{
		Address: host,
		Port:    self.listenPort,
		ID:      self.ID(),
		Name:    self.name,
	}

	return discovery.Default.Register(sd)

}

func (self *cellService) Run() error {

	err := self.Start()
	if err != nil {
		return err
	}

	util.WaitExit()

	return self.Stop()
}
func (self *cellService) Stop() error {
	return discovery.Default.Deregister(self.ID())
}

func NewService(name string) Service {

	return &cellService{
		name: name,
	}
}
