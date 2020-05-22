package redsd

import (
	"sync"
)

type NodeList struct {
	name      string
	nodeByID  map[string]*NodeContext
	nodeList  []*NodeDesc
	dirty     bool
	guard     sync.Mutex
	sd        *RedisDiscovery
	kind      int
	heartBeat bool
}

func (self *NodeList) Kind() int {
	return self.kind
}

func (self *NodeList) Name() string {
	return self.name
}

func (self *NodeList) DescList() (ret []*NodeDesc) {

	self.guard.Lock()
	defer self.guard.Unlock()

	if self.dirty {

		ret = make([]*NodeDesc, 0, len(self.nodeByID))
		for _, ctx := range self.nodeByID {
			ret = append(ret, ctx.Desc)
		}

		self.nodeList = ret

		self.dirty = false
	} else {
		ret = self.nodeList
	}

	return
}

func (self *NodeList) AddDesc(ctx *NodeContext) {
	self.guard.Lock()
	self.nodeByID[ctx.Desc.ID] = ctx
	self.dirty = true
	self.guard.Unlock()
}

func (self *NodeList) DeleteDesc(nodeid string) {
	self.guard.Lock()
	self.dirty = true
	delete(self.nodeByID, nodeid)
	self.guard.Unlock()
}

func (self *NodeList) GetDesc(nodeid string) *NodeContext {
	self.guard.Lock()
	defer self.guard.Unlock()

	if ctx, ok := self.nodeByID[nodeid]; ok {
		return ctx
	}

	return nil
}

func newNodeList(name string, sd *RedisDiscovery) *NodeList {
	return &NodeList{
		name:     name,
		dirty:    true,
		nodeByID: make(map[string]*NodeContext),
		sd:       sd,
	}
}
