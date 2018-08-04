package cellep

import (
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	_ "github.com/davyxu/cellmesh/discovery/consul"
	"github.com/davyxu/cellmesh/endpoint"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
)

type cellEndPoint struct {
	name       string
	listenPort int
	svcByName  sync.Map // map[reflect.Type]*endpoint.ServiceInfo
}

func (self *cellEndPoint) AddHandler(name string, svc *endpoint.ServiceInfo) {

	self.svcByName.Store(svc.RequestType, svc)
}

func (self *cellEndPoint) ID() string {
	return fmt.Sprintf("%s-%d", self.name, self.listenPort)
}

func (self *cellEndPoint) Start() error {

	p := peer.NewGenericPeer("tcp.Acceptor", "node", ":0", nil)

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {

		msgType := reflect.TypeOf(ev.Message()).Elem()

		if svcRaw, ok := self.svcByName.Load(msgType); ok {

			svc := svcRaw.(*endpoint.ServiceInfo)

			e := &endpoint.Event{
				Request:  ev.Message(),
				Response: svc.NewResponse(),
			}

			svc.Handler(e)

			ev.Session().Send(e.Response)
		}

	})

	p.Start()

	self.listenPort = p.(cellnet.TCPAcceptor).ListenPort()

	host := util.GetLocalIP()

	sd := &discovery.ServiceDesc{
		Address: host,
		Port:    self.listenPort,
		ID:      self.ID(),
		Name:    self.name,
	}

	return discovery.Default.Register(sd)

}

func (self *cellEndPoint) Run() error {

	err := self.Start()
	if err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch

	return self.Stop()
}
func (self *cellEndPoint) Stop() error {

	//return discovery.Default.Deregister(self.ID())
	return nil
}

func NewService(name string) endpoint.EndPoint {

	return &cellEndPoint{
		name: name,
	}
}
