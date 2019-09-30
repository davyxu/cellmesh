package memsd

import (
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/svc/memsd/model"
	"github.com/davyxu/cellmesh/svc/memsd/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"time"
)

func (self *memDiscovery) clearCache() {

	self.svcCacheGuard.Lock()
	self.svcCache = map[string][]*discovery.ServiceDesc{}
	self.svcCacheGuard.Unlock()

	self.kvCacheGuard.Lock()
	self.kvCache = map[string][]byte{}
	self.kvCacheGuard.Unlock()
}

func (self *memDiscovery) connect(addr string) {
	p := peer.NewGenericPeer("tcp.Connector", "memsd", addr, self.q)

	proc.BindProcessorHandler(p, "memsd.cli", func(ev cellnet.Event) {

		switch msg := ev.Message().(type) {
		case *cellnet.SessionConnected:

			self.sesGuard.Lock()
			self.ses = ev.Session()
			self.sesGuard.Unlock()
			self.clearCache()
			ev.Session().Send(&sdproto.AuthREQ{
				Token: self.token,
			})
		case *cellnet.SessionClosed:
			self.token = ""
			log.Errorf("memsd discovery lost!")

		case *sdproto.AuthACK:

			self.token = msg.Token

			if self.initWg != nil {
				// Pull的消息还要在queue里处理，这里确认处理完成后才算初始化完成
				self.initWg.Done()
			}

			log.SetColor("green").Debugf("memsd discovery ready!")

			self.triggerNotify("ready")

		case *sdproto.ValueChangeNotifyACK:

			if model.IsServiceKey(msg.Key) {
				self.updateSvcCache(msg.SvcName, msg.Value)
			} else {
				self.updateKVCache(msg.Key, msg.Value)
			}

		case *sdproto.ValueDeleteNotifyACK:

			if model.IsServiceKey(msg.Key) {
				svcid := model.GetSvcIDByServiceKey(msg.Key)
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

	go func() {

		ticker := time.NewTicker(time.Second * 5)
		for {
			<-ticker.C

			self.ses.Send(&sdproto.SDPingACK{})
		}

	}()

}
