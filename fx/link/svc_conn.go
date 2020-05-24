package link

import (
	"github.com/davyxu/cellmesh/fx/redsd"
	"github.com/davyxu/cellmesh/svc/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/ulog"
	"time"
)

// 连接到服务
func ConnectNode(param *NodeParameter) {

	nodeList := SD.NewNodeList(param.NodeName, int(proto.NodeKind_Connect))

	// 实时检查更新
	go func() {
		for {

			addList, delList := nodeList.Check()

			for _, ctx := range delList {

				nodeList.DeleteDesc(ctx.Desc.ID)

				ulog.WithField("nodeid", ctx.Desc.ID).Debugf("discovery conn delete link")
				closeNode(ctx.Desc)
			}

			for _, ctx := range addList {
				ulog.WithField("nodeid", ctx.Desc.ID).Debugf("discovery conn add link")
				nodeList.AddDesc(ctx)
				connNode(ctx.Desc, param)
			}

			time.Sleep(time.Second * 3)
		}
	}()

}

func connNode(desc *redsd.NodeDesc, param *NodeParameter) {

	p := peer.NewGenericPeer(param.PeerType, param.NodeName, desc.Address(), param.Queue)

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

func closeNode(desc *redsd.NodeDesc) {
	if desc.Peer != nil {
		desc.Peer.Stop()
		desc.Peer = nil
		desc.Session = nil
	}
}
