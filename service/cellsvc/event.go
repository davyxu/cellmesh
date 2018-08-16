package cellsvc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/relay"
)

type svcEvent struct {
	cellnet.Event
}

type ReplyEvent interface {
	Reply(msg interface{})
}

func (self *svcEvent) GetContextID() int64 {
	if relayEvent, ok := self.Event.(*relay.RecvMsgEvent); ok {
		return relayEvent.OneContextID()
	}

	return 0
}

func (self *svcEvent) Reply(msg interface{}) {

	if replyEv, ok := self.Event.(ReplyEvent); ok {
		replyEv.Reply(msg)
	} else {
		panic("Require 'ReplyEvent' to reply event")
	}
}
