package cellsvc

import (
	"fmt"
	"github.com/davyxu/cellmicro"
	"github.com/davyxu/cellmicro/discovery"
	_ "github.com/davyxu/cellmicro/discovery/consul"
	"github.com/davyxu/cellmicro/svc"
	"github.com/davyxu/cellmicro/util"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"os"
	"os/signal"
	"syscall"
)

type cellService struct {
	name          string
	port          int
	handlerByName map[string]svc.Handler
}

func (self *cellService) AddHandler(name string, handler svc.Handler) {

	self.handlerByName[name] = handler
}

func (self *cellService) ID() string {
	return fmt.Sprintf("%s-%d", self.name, self.port)
}

func (self *cellService) Start() error {

	q := cellnet.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Acceptor", "node", fmt.Sprintf(":%d", self.port), q)

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {

		meta := cellnet.MessageMetaByMsg(ev.Message())
		if meta != nil {

			msgName := meta.FullName()

			respMeta := cellmicro.GetResponseMeta(meta)

			if h, ok := self.handlerByName[msgName]; ok {

				e := &svc.Event{
					Request:  ev.Message(),
					Response: respMeta.NewType(),
				}

				h(e)

				ev.Session().Send(e.Response)
			}

		}
	})

	p.Start()

	q.StartLoop()

	addrss := util.GetLocalIP()

	sd := &discovery.ServiceDesc{
		Address: addrss,
		Port:    self.port,
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

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch

	return self.Stop()
}
func (self *cellService) Stop() error {

	return discovery.Default.Deregister(self.ID())
}

func (self *cellService) SetPort(port int) {
	self.port = port
}

func NewService(name string) svc.Service {

	return &cellService{
		name:          name,
		port:          14330,
		handlerByName: make(map[string]svc.Handler),
	}
}
