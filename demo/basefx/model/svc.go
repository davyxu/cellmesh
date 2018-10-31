package fxmodel

import "github.com/davyxu/cellnet"

var (
	Queue cellnet.EventQueue
)

type ServiceParameter struct {
	SvcName      string // 服务名,注册到服务发现
	NetProcName  string // cellnet处理器名称
	NetPeerType  string // cellnet的PeerType
	ListenAddr   string // socket侦听地址
	MaxConnCount int    // 最大连接数量
	NoQueue      bool   // 不使用队列
}
