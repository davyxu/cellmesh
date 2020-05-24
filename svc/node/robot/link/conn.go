package link

import (
	"fmt"
	"github.com/davyxu/cellmesh/svc/node/robot/model"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"strings"
)

func ConnectTCP(r *model.Robot, peerName, address string) {

	r.SetState("Connect" + strings.Title(peerName))
	p := peer.NewGenericPeer("tcp.SyncConnector", fmt.Sprintf("%s%s", peerName, r.ID), address, nil)
	r.SetPeer(peerName, p)

	proc.BindProcessorHandler(p, "tcp.robot", nil, r)

	p.Start()
}
