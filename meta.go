package cellmicro

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/util"
)

var (
	requestMapper = map[*cellnet.MessageMeta]*cellnet.MessageMeta{}
)

func RegisterRequestPair(req, ack *cellnet.MessageMeta) {

	req.Codec = codec.MustGetCodec("json")
	req.ID = int(util.StringHash(req.FullName()))
	ack.Codec = codec.MustGetCodec("json")
	ack.ID = int(util.StringHash(ack.FullName()))

	cellnet.RegisterMessageMeta(req)
	cellnet.RegisterMessageMeta(ack)

	requestMapper[req] = ack
}

func GetResponseMeta(req *cellnet.MessageMeta) *cellnet.MessageMeta {
	return requestMapper[req]
}
