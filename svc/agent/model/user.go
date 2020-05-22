package model

import (
	"github.com/davyxu/cellmesh/link"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
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

	CID proto.ClientID
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

func (self *User) SendToBackend(backendSvcid string, msgID int, msgData []byte) {

	logfields := ulog.Fields{
		"sesid":   self.Session.ID(),
		"nodeid":  backendSvcid,
		"msgid":   msgID,
		"msgsize": len(msgData),
	}

	desc := link.DescByID(backendSvcid)

	if desc == nil {
		ulog.WithFields(logfields).Errorf("backend node not found")
		return
	}

	backendSes := link.LinkByDesc(desc)

	if backendSes == nil {
		ulog.WithFields(logfields).Errorf("backend node not ready")
		return
	}

	userMsg, meta, _ := codec.DecodeMessage(int(msgID), msgData)

	if meta != nil {
		logfields["msgname"] = meta.FullName()
		msgtostr := cellnet.MessageToString(userMsg)
		logfields["msgtostr"] = msgtostr

		//ulog.Debugf("client(%d) -> %s len:%d %s| %s", self.Session.ID(), backendSvcid, len(msgData), meta.FullName(), msgtostr)
	}

	backendSes.Send(&proto.RouterTransmitACK{
		MsgID:    uint32(msgID),
		MsgData:  msgData,
		ClientID: self.CID.ID,
	})
}

// 绑定用户后台
func (self *User) BindBackend(svcName string, nodeid string) {

	for _, t := range self.Targets {
		if t.Name == svcName {
			t.NodeID = nodeid
			return
		}
	}

	self.CID = proto.ClientID{
		ID:    self.Session.ID(),
		SvcID: AgentSvcID,
	}

	self.Targets = append(self.Targets, &Backend{
		Name:   svcName,
		NodeID: nodeid,
	})

	ulog.WithFields(ulog.Fields{
		"sesid":  self.Session.ID(),
		"nodeid": nodeid,
	}).Debugf("user bind backend")
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
