package redsd

import (
	"github.com/davyxu/ulog"
	"github.com/gomodule/redigo/redis"
	"time"
)

func (self *NodeList) Monitor(interval time.Duration) {
	// 实时检查更新
	go func() {
		for {

			addList, delList := self.Check()

			for _, ctx := range delList {
				ulog.WithField("nodeid", ctx.Desc.ID).Debugf("discovery monitor delete link")
				self.DeleteDesc(ctx.Desc.ID)
			}

			for _, ctx := range addList {
				ulog.WithField("nodeid", ctx.Desc.ID).Debugf("discovery monitor add link")
				self.AddDesc(ctx)
			}

			time.Sleep(interval)
		}
	}()

}

// 刷新关心的节点信息
func (self *NodeList) Check() (addList, delList []*NodeContext) {

	self.sd.operate(func(conn redis.Conn) {

		nodeIDList, err := redis.Strings(conn.Do("HKEYS", key_NodeReg+self.name))
		if err != nil {
			ulog.Errorf("node reg fetch failed: %s", err)
			return
		}

		// 遍历每种服务的每一个ID
		for _, nodeID := range nodeIDList {

			nodeKey := key_Node + nodeID

			ctx := self.GetDesc(nodeID)

			ver, err := redis.Int64(conn.Do("HGET", nodeKey, key_NodeVer))

			if err != nil {

				if ctx != nil {
					delList = append(delList, ctx)
				}

				continue
			}

			var create bool

			if ctx != nil {

				if ctx.Ver == ver {
					continue
				}

			} else {

				ctx = &NodeContext{
					Ver: ver,
				}

				create = true
			}

			// 维护redis节点正确性
			data, err := redis.Bytes(conn.Do("HGET", nodeKey, key_NodeDesc))
			if err != nil {
				continue
			}

			if create {
				addList = append(addList, ctx)
			} else {
				delList = append(delList, ctx)
			}

			ctx.Desc = unmarshalDesc(data)

		}
	})

	return
}
