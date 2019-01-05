package memsd

import (
	"github.com/davyxu/cellmesh/discovery/memsd/model"
	"github.com/davyxu/cellmesh/discovery/memsd/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"strings"
	"time"
)

func (self *memDiscovery) connect(addr string) {
	p := peer.NewGenericPeer("tcp.Connector", "memsd", addr, model.Queue)

	proc.BindProcessorHandler(p, "memsd.cli", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *cellnet.SessionConnected:
			self.sesGuard.Lock()
			self.ses = p.(cellnet.TCPConnector).Session()
			self.sesGuard.Unlock()
		case *proto.ValueChangeNotifyACK:

			if strings.HasPrefix(msg.Key, servicePrefix) {
				self.updateSvcCache(msg.SvcName, msg.Value)
			} else {
				self.updateKVCache(msg.Key, msg.Value)
			}

		case *proto.ValueDeleteNotifyACK:

			if strings.HasPrefix(msg.Key, servicePrefix) {
				svcid := msg.Key[len(servicePrefix):]
				self.deleteSvcCache(svcid, msg.SvcName)
			} else {
				self.deleteKVCache(msg.Key)
			}
		}
	})

	// noDelay
	p.(cellnet.TCPSocketOption).SetSocketBuffer(1024*1024, 1024*1024, true)

	// 断线后自动重连
	p.(cellnet.TCPConnector).SetReconnectDuration(time.Second * 5)

	p.Start()

	for {

		if p.(cellnet.PeerReadyChecker).IsReady() {
			break
		}

		time.Sleep(time.Millisecond * 500)
	}

}
