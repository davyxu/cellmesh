package model

import (
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"strings"
)

const (
	ServiceKeyPrefix = "_svcdesc_"
)

var (
	Queue cellnet.EventQueue
	IDGen = meshutil.NewUUID64Generator()

	Listener cellnet.Peer
	Debug    bool

	Version = "0.1.0"
)

func IsServiceKey(rawkey string) bool {

	return strings.HasPrefix(rawkey, ServiceKeyPrefix)
}

func GetSvcIDByServiceKey(rawkey string) string {

	if IsServiceKey(rawkey) {
		return rawkey[len(ServiceKeyPrefix):]
	}

	return ""
}

func init() {
	IDGen.AddTimeComponent(8)
	IDGen.AddSeqComponent(8, 0)
}

func GetSessionToken(ses cellnet.Session) (token string) {
	ses.(cellnet.ContextSet).FetchContext("token", &token)

	return
}

func Broadcast(msg interface{}) {
	Listener.(cellnet.TCPAcceptor).VisitSession(func(ses cellnet.Session) bool {
		ses.Send(msg)
		return true
	})
}

func TokenExists(token string) (ret bool) {
	Listener.(cellnet.TCPAcceptor).VisitSession(func(ses cellnet.Session) bool {

		if GetSessionToken(ses) == token {
			ret = true
			return false
		}

		return true
	})

	return
}
