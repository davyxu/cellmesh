package agentapi

import (
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/ulog"
)

type AgentMsgEvent struct {
	Ses      cellnet.Session
	Msg      interface{}
	ClientID int64
}

func (self *AgentMsgEvent) Session() cellnet.Session {
	return self.Ses
}

func (self *AgentMsgEvent) Message() interface{} {
	return self.Msg
}

func (self *AgentMsgEvent) Reply(msg interface{}) {

	data, meta, err := codec.EncodeMessage(msg, nil)
	if err != nil {
		ulog.Errorf("Reply.EncodeMessage %s", err)
		return
	}

	self.Ses.Send(&proto.RouterTransmitACK{
		MsgID:    uint32(meta.ID),
		MsgData:  data,
		ClientID: self.ClientID,
	})

}
