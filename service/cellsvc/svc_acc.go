package cellsvc

import (
	"fmt"
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
	evDispatcher

	svcName    string
	listenAddr string
	listener   cellnet.GenericPeer
}

func (self *accService) String() string {
	return fmt.Sprintf("Acceptor: '%s' addr: '%s'", self.svcName, self.listenAddr)
}

func (self *accService) IsReady() bool {
	if self.listener == nil {
		return false
	}

	return self.listener.(interface {
		IsReady() bool
	}).IsReady()
}

func (self *accService) Start() {

	self.listener = peer.NewGenericPeer("tcp.Acceptor", self.svcName, self.listenAddr, nil)

	proc.BindProcessorHandler(self.listener, self.procName, func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *proto.ServiceIdentifyACK:

			if pre := service.GetRemoteService(msg.SvcID); pre == nil {

				// 添加连接上来的对方服务
				service.AddRemoteService(ev.Session(), &discovery.ServiceDesc{
					ID:   msg.SvcID,
					Name: msg.SvcName,
				})
			}

		case *cellnet.SessionClosed:
			service.RemoveRemoteService(ev.Session())
		}

		self.Invoke(&svcEvent{
			Event: ev,
		})
	})

	self.listener.Start()

	host := util.GetLocalIP()

	sd := &discovery.ServiceDesc{
		Host: host,
		Port: self.listener.(cellnet.TCPAcceptor).Port(),
		ID:   fxmodel.GetSvcID(self.svcName),
		Name: self.svcName,
		Tags: []string{fxmodel.Node},
	}

	log.SetColor("green").Debugf("service '%s' listen at %s:%d", sd.ID, host, sd.Port)

	// TODO 注册之前先搜索，有重名的启动失败
	discovery.Default.Register(sd)
}

func (self *accService) Stop() {
	discovery.Default.Deregister(fxmodel.GetSvcID(self.svcName))
}

// listenAddr 格式:
// :0，自动设置端口，
// host:min~max设置[min,max]范围的可用端口
func NewCommunicateAcceptor(svcName, listenAddr string) service.CommunicateService {

	if listenAddr == "" {
		listenAddr = ":0"
	}

	return &accService{
		listenAddr: listenAddr,
		svcName:    svcName,
	}
}
