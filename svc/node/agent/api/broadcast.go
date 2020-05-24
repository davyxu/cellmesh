package agentapi

import (
	"github.com/davyxu/cellmesh/fx/link"
	"github.com/davyxu/cellmesh/svc/proto"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/ulog"
)

// 关闭所有网关上客户端的连接
func CloseAllClient() {

	for _, desc := range link.DescListByName("backend") {
		ses := link.LinkByDesc(desc)
		if ses != nil {
			ses.Send(&proto.AgentCloseClientACK{
				All: true,
			})
		}
	}
}

// 广播给所有客户端
func BroadcastAll(msg interface{}) {

	data, meta, err := codec.EncodeMessage(msg, nil)
	if err != nil {
		ulog.Errorf("BroadcastAll.EncodeMessage %s", err)
		return
	}

	for _, desc := range link.DescListByName("backend") {
		ses := link.LinkByDesc(desc)
		if ses != nil {
			ses.Send(&proto.AgentTransmitACK{
				MsgID:   uint32(meta.ID),
				MsgData: data,
				All:     true,
			})
		}
	}

}

// 给客户端发消息
func Send(cid *proto.AgentClientID, msg interface{}) {

	agentSes := link.LinkByID(cid.NodeID)
	if agentSes != nil {
		data, meta, err := codec.EncodeMessage(msg, nil)
		if err != nil {
			ulog.Errorf("Send.EncodeMessage %s", err)
			return
		}

		agentSes.Send(&proto.AgentTransmitACK{
			MsgID:    uint32(meta.ID),
			MsgData:  data,
			ClientID: cid.SessionID,
		})
	}
}

type ClientList struct {
	sesByAgentSvcID map[string][]int64
}

// 添加客户端
func (self *ClientList) AddClient(cid *proto.AgentClientID) {
	seslist := self.sesByAgentSvcID[cid.NodeID]
	seslist = append(seslist, cid.SessionID)
	self.sesByAgentSvcID[cid.NodeID] = seslist
}

// 关闭列表中客户端的连接
func (self *ClientList) CloseClient() {
	for agentSvcID, sesList := range self.sesByAgentSvcID {

		agentSes := link.LinkByID(agentSvcID)
		if agentSes != nil {
			agentSes.Send(&proto.AgentCloseClientACK{
				ID: sesList,
			})
		}
	}
}

// 将消息广播给列表中的客户端
func (self *ClientList) Broadcast(msg interface{}) {

	data, meta, err := codec.EncodeMessage(msg, nil)
	if err != nil {
		ulog.Errorf("ClientList.EncodeMessage %s", err)
		return
	}

	for agentSvcID, sesList := range self.sesByAgentSvcID {

		agentSes := link.LinkByID(agentSvcID)
		if agentSes != nil {

			agentSes.Send(&proto.AgentTransmitACK{
				MsgID:        uint32(meta.ID),
				MsgData:      data,
				ClientIDList: sesList,
			})

		} else {
			ulog.Warnf("Agent not ready, ignore msg, svcid: '%s' msg: '%+v'", agentSvcID, msg)
		}
	}
}

// 创建一个客户端列表
func NewClientList() *ClientList {

	return &ClientList{
		sesByAgentSvcID: make(map[string][]int64),
	}
}
