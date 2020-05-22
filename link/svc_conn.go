package link

import (
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellmesh/redsd"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/ulog"
	"time"
)

// 连接到服务
func ConnectNode(param *NodeParameter) {

	nodeList := SD.NewNodeList(param.SvcName, int(proto.NodeKind_Connect))

	// 实时检查更新
	go func() {
		for {

			addList, delList := nodeList.Check()

			for _, ctx := range delList {

				nodeList.DeleteDesc(ctx.Desc.ID)

				ulog.WithField("nodeid", ctx.Desc.ID).Debugf("discovery delete link")
				if ctx.Desc.Peer != nil {
					ctx.Desc.Peer.Stop()
					ctx.Desc.Peer = nil
					ctx.Desc.Session = nil
				}
			}

			for _, ctx := range addList {
				ulog.WithField("nodeid", ctx.Desc.ID).Debugf("discovery add link")
				nodeList.AddDesc(ctx)
				connNode(ctx.Desc, param)
			}

			time.Sleep(time.Second * 3)
		}
	}()

}

func connNode(desc *redsd.NodeDesc, param *NodeParameter) {

	p := peer.NewGenericPeer(param.PeerType, param.SvcName, desc.Address(), param.Queue)

	proc.BindProcessorHandler(p, param.NetProc, param.EventCallback)

	if opt, ok := p.(cellnet.TCPSocketOption); ok {
		opt.SetSocketBuffer(2048, 2048, true)
	}

	p.(cellnet.TCPConnector).SetReconnectDuration(time.Second * 5)

	p.Start()

	desc.Peer = p

	ctxSet := p.(cellnet.ContextSet)
	ctxSet.SetContext("NodeDesc", desc)
}
