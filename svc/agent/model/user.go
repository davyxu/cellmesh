package model

import (
	"errors"
	"github.com/davyxu/cellmesh/link"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
	"time"
)

type Backend struct {
	SvcName string
	SvcID   string // 只保留绑定后台的svcid,即便后台更换session,也无需同步
}

type User struct {
	ClientSession cellnet.Session
	Targets       []*Backend
	LastPingTime  time.Time

	CID proto.ClientID
}

// 广播到这个用户绑定的所有后台
func (self *User) BroadcastToBackends(msg interface{}) {

	for _, t := range self.Targets {

		backendSes := link.GetLink(t.SvcID)
		if backendSes != nil {
			backendSes.Send(msg)
		}
	}
}

var (
	ErrBackendNotFound = errors.New("backend not found")
)

func (self *User) TransmitToBackend(backendSvcid string, msgID int, msgData []byte) error {

	backendSes := link.GetLink(backendSvcid)

	if backendSes == nil {
		return ErrBackendNotFound
	}

	backendSes.Send(&proto.RouterTransmitACK{
		MsgID:    uint32(msgID),
		MsgData:  msgData,
		ClientID: self.CID.ID,
	})

	return nil
}

// 绑定用户后台
func (self *User) SetBackend(svcName string, svcID string) {

	for _, t := range self.Targets {
		if t.SvcName == svcName {
			t.SvcID = svcID
			return
		}
	}

	self.CID = proto.ClientID{
		ID:    self.ClientSession.ID(),
		SvcID: AgentSvcID,
	}

	self.Targets = append(self.Targets, &Backend{
		SvcName: svcName,
		SvcID:   svcID,
	})
}

func (self *User) GetBackend(svcName string) string {

	for _, t := range self.Targets {
		if t.SvcName == svcName {
			return t.SvcID
		}
	}

	return ""
}

func NewUser(clientSes cellnet.Session) *User {
	return &User{
		ClientSession: clientSes,
	}
}
