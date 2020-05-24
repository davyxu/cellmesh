package model

type FrontendParameter struct {
	SvcName     string // 服务名,注册到服务发现
	NetProcName string // cellnet处理器名称
	NetPeerType string // cellnet的PeerType
	ListenAddr  string // socket侦听地址
}
