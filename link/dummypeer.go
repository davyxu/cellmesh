package link

import "github.com/davyxu/cellnet"

type dummyPeer struct {
	*ServiceParameter
}

func (self *dummyPeer) IsDummy() bool {
	return true
}

// 获取SetAddress中的侦听或者连接地址
func (self *dummyPeer) Address() string {

	return ""
}

func (self *dummyPeer) Start() cellnet.Peer {
	return self
}

func (self *dummyPeer) Stop() {

}

func (self *dummyPeer) Name() string {
	return self.SvcName
}

func (self *dummyPeer) IsReady() bool {
	return false
}

func (self *dummyPeer) TypeName() string {
	return self.PeerType
}

func (self *dummyPeer) Session() cellnet.Session {
	return nil
}

func newDummyPeer(param *ServiceParameter) cellnet.Peer {
	return &dummyPeer{
		ServiceParameter: param,
	}
}
