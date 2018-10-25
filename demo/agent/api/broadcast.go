package api

import (
	"github.com/davyxu/cellmesh/demo/agent/model"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/relay"
)

// 关闭所有网关上客户端的连接
func CloseAllClient() {

	service.VisitRemoteService(func(ses cellnet.Session, ctx *service.RemoteServiceContext) bool {

		if ctx.Name == model.BackendName {
			ses.Send(&proto.CloseClientACK{
				All: true,
			})
		}

		return true
	})
}

// 广播给所有客户端
func BroadcastAll(msg interface{}) {
	service.VisitRemoteService(func(ses cellnet.Session, ctx *service.RemoteServiceContext) bool {

		if ctx.Name == model.BackendName {
			relay.Relay(ses, msg)
		}

		return true
	})
}

type ClientList struct {
	sesByAgentSvcID map[string][]int64
}

// 添加客户端
func (self *ClientList) AddClient(cid proto.ClientID) {
	seslist := self.sesByAgentSvcID[cid.SvcID]
	seslist = append(seslist, cid.ID)
	self.sesByAgentSvcID[cid.SvcID] = seslist
}

// 关闭列表中客户端的连接
func (self *ClientList) CloseClient() {
	for agentSvcID, sesList := range self.sesByAgentSvcID {

		agentSes := service.GetRemoteService(agentSvcID)
		if agentSes != nil {
			agentSes.Send(proto.CloseClientACK{
				ID: sesList,
			})
		}
	}
}

// 将消息广播给列表中的客户端
func (self *ClientList) Broadcast(msg interface{}) {

	for agentSvcID, sesList := range self.sesByAgentSvcID {

		agentSes := service.GetRemoteService(agentSvcID)
		if agentSes != nil {

			relay.Relay(agentSes, msg, sesList)
		} else {
			log.Warnf("Agent not ready, ignore msg, svcid: '%s' msg: '%+v'", agentSvcID, msg)
		}
	}
}

// 创建一个客户端列表
func NewClientList() *ClientList {

	return &ClientList{
		sesByAgentSvcID: make(map[string][]int64),
	}
}
