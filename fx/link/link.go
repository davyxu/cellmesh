package link

import (
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/fx/redsd"
	"github.com/davyxu/cellmesh/svc/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/ulog"
)

// 某类节点的默认连接
func LinkByName(name string) cellnet.Session {

	descList := DescListByName(name)
	if len(descList) == 0 {
		return nil
	}

	return LinkByDesc(descList[0])
}

// 根据ID获取会话
func LinkByID(nodeid string) (ret cellnet.Session) {

	desc := DescByID(nodeid)
	if desc == nil {
		return nil
	}

	return LinkByDesc(desc)
}

func LinkByDesc(desc *redsd.NodeDesc) cellnet.Session {

	// Acceptor 连接上来的连接
	if desc.Session != nil {
		return desc.Session
	}

	if desc.Peer == nil {
		return nil
	}

	return desc.Peer.(interface {
		// 默认会话
		Session() cellnet.Session
	}).Session()
}

// 某类节点的list
func DescListByName(name string) []*redsd.NodeDesc {
	nodeList := SD.NodeListByName(name)
	if nodeList == nil {
		return nil
	}

	return nodeList.DescList()
}

// 某节点的desc
func DescByID(nodeid string) *redsd.NodeDesc {
	name, _, _, err := fx.ParseNodeID(nodeid)
	if err != nil {
		return nil
	}

	nodeList := SD.NodeListByName(name)
	if nodeList == nil {
		return nil
	}

	ctx := nodeList.GetDesc(nodeid)

	if ctx == nil {
		return nil
	}

	return ctx.Desc
}

func DescByLink(ses cellnet.Session) *redsd.NodeDesc {
	if ses == nil {
		return nil
	}

	if raw, ok := ses.(cellnet.ContextSet).GetContext("NodeDesc"); ok {
		return raw.(*redsd.NodeDesc)
	}

	return nil
}

func addLink(desc *redsd.NodeDesc) {
	nodeList := SD.NodeListByName(desc.Name)
	if nodeList != nil && nodeList.GetDesc(desc.ID) != nil {
		ulog.Warnf("duplicate node: %s", desc.ID)
	}

	if nodeList == nil {
		nodeList = SD.NewNodeList(desc.Name, int(proto.NodeKind_Accept))
	}

	nodeList.AddDesc(&redsd.NodeContext{
		Desc: desc,
		Ver:  1,
	})
}

func removeLink(desc *redsd.NodeDesc) {
	nodeList := SD.NodeListByName(desc.Name)
	if nodeList == nil {
		return
	}

	nodeList.DeleteDesc(desc.ID)
}
