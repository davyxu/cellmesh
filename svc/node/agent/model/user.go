package model

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/fx/link"
	"github.com/davyxu/cellmesh/svc/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/ulog"
	"time"
)

type Backend struct {
	Name   string
	NodeID string // 只保留绑定后台的svcid,即便后台更换session,也无需同步
}

type User struct {
	Session      cellnet.Session
	Targets      []*Backend
	LastPingTime time.Time

	CID proto.AgentClientID
}

// 广播到这个用户绑定的所有后台
func (self *User) BroadcastToBackends(msg interface{}) {

	for _, t := range self.Targets {

		backendSes := link.LinkByID(t.NodeID)
		if backendSes != nil {
			backendSes.Send(msg)
		}
	}
}

func (self *User) SendToBackend(nodeID string, msgID int, msgData []byte) {
	SendToBackend(nodeID, msgID, msgData, self.CID.SessionID)
}

func SendToBackend(nodeID string, msgID int, msgData []byte, clientSesID int64) {

	desc := link.DescByID(nodeID)

	if desc == nil {
		ulog.Errorf("backend node not found, nodeid: %s msgid: %d", nodeID, msgID)
		return
	}

	backendSes := link.LinkByDesc(desc)

	if backendSes == nil {
		ulog.Errorf("backend node not ready, nodeid: %s msgid: %d", nodeID, msgID)
		return
	}

	backendSes.Send(&proto.AgentTransmitACK{
		MsgID:    uint32(msgID),
		MsgData:  msgData,
		ClientID: clientSesID,
	})
}

// 绑定用户后台
func (self *User) BindBackend(nodeid string) {

	nodeName, _, _, _ := fx.ParseNodeID(nodeid)

	for _, t := range self.Targets {
		if t.Name == nodeName {
			t.NodeID = nodeid
			return
		}
	}

	self.CID = proto.AgentClientID{
		SessionID: self.Session.ID(),
		NodeID:    AgentNodeID,
	}

	self.Targets = append(self.Targets, &Backend{
		Name:   nodeName,
		NodeID: nodeid,
	})

	ulog.Debugf("user bind backend, sesid: %d nodeid: %s", self.Session.ID(), nodeid)
}

func (self *User) GetBackend(svcName string) string {

	for _, t := range self.Targets {
		if t.Name == svcName {
			return t.NodeID
		}
	}

	return ""
}

func NewUser(clientSes cellnet.Session) *User {
	return &User{
		Session: clientSes,
	}
}
