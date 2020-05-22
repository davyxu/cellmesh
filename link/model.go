package link

import "github.com/davyxu/cellmesh/redsd"

var (
	descBySvcID = map[string]*redsd.NodeDesc{}
	SD          *redsd.RedisDiscovery
)
