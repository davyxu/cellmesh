package model

import (
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/relay"
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

		backendSes := service.GetRemoteService(t.SvcID)
		if backendSes != nil {
			backendSes.Send(msg)
		}
	}
}

// 将消息发到绑定的某类服务器
func (self *User) RelayToService(svcName string, msg interface{}) {
	backendSvcid := self.GetBackend(svcName)

	backendSes := service.GetRemoteService(backendSvcid)

	if backendSes != nil {
		relay.Relay(backendSes, msg, self.CID.ID, self.CID.SvcID)
	} else {
		log.Warnf("Backend not found, msg: '%s' mode: 'auth'", cellnet.MessageToName(msg))
	}
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
