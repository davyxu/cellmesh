package cellmesh

import (
	"github.com/davyxu/cellnet"
)

// 回复event来源一个消息
func Reply(ev cellnet.Event, msg interface{}) {

	type replyEvent interface {
		Reply(msg interface{})
	}

	if replyEv, ok := ev.(replyEvent); ok {
		replyEv.Reply(msg)
	} else {
		panic("Require 'ReplyEvent' to reply event")
	}
}
