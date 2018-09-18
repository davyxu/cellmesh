package cellsvc

import (
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx/model"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"strings"
	"sync"
	"time"
)

type connector interface {
	cellnet.TCPConnector
	IsReady() bool
}

type conService struct {
	evDispatcher

	svcName string // 让远程服务看到的自己的服务名

	tgtSvcName string // 远程要连接的服务名

	connectorBySvcID sync.Map // map[svcid] connector
}

func (self *conService) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Connector '%s' ", self.tgtSvcName))

	self.connectorBySvcID.Range(func(key, value interface{}) bool {

		conn := value.(connector)

		//var sd *discovery.ServiceDesc
		//conn.(cellnet.ContextSet).FetchContext("sd", &sd)
		pp := conn.(cellnet.PeerProperty)

		sb.WriteString(fmt.Sprintf("addr: '%s' ready: %v", pp.Address(), conn.IsReady()))

		return true
	})

	return sb.String()
}

func (self *conService) IsReady() (ret bool) {

	ret = true

	var count int
	self.connectorBySvcID.Range(func(key, value interface{}) bool {

		count++

		conn := value.(connector)

		if !conn.IsReady() {
			ret = false
			return false
		}

		return true
	})

	return count > 0 && ret
}

func (self *conService) connFlow(p cellnet.GenericPeer, sd *discovery.ServiceDesc) {

	var stop sync.WaitGroup

	p.(cellnet.ContextSet).SetContext("sd", sd)

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

		descList, err := discovery.Default.Query(self.tgtSvcName)

		descList = discovery.MatchAnyTag(descList, fxmodel.MatchNodes...)

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

// 连接目标服务,并告诉对方自己服务名字
func NewCommunicateConnector(svcName, tgtSvcName string) service.CommunicateService {

	return &conService{
		tgtSvcName: tgtSvcName,
		svcName:    svcName,
	}
}
