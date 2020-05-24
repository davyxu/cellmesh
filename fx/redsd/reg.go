package redsd

import (
	"github.com/davyxu/ulog"
	"github.com/gomodule/redigo/redis"
)

// 清理本类型的未使用的节点
func (self *NodeList) recycle(name string) {
	self.sd.operate(func(conn redis.Conn) {

		regKey := key_NodeReg + name
		nodeIDList, err := redis.Strings(conn.Do("HKEYS", regKey))
		if err != nil {
			ulog.Errorf("node reg fetch failed: %s", err)
			return
		}

		// 遍历每种服务的每一个ID
		for _, nodeID := range nodeIDList {

			nodeKey := key_Node + nodeID

			safeOpNode(conn, nodeID, func() {

				_, err := redis.Int64(conn.Do("HGET", nodeKey, key_NodeVer))
				if err != nil {
					ulog.WithField("nodeid", nodeID).Debugln("remove unused node")
					// 删除已经不存在的node注册信息
					exec(conn, "HDEL", regKey, nodeID)
				}
			})
		}
	})
}

func (self *NodeList) Register(desc *NodeDesc) {
	self.heartBeat = true

	self.recycle(self.name)

	self.sd.operate(func(conn redis.Conn) {

		safeOpNode(conn, desc.ID, func() {
			// 注册入口
			exec(conn, "HSET", key_NodeReg+desc.Name, desc.ID, 1)

			nodeKey := key_Node + desc.ID
			// 本体数据
			exec(conn, "HSET", nodeKey, key_NodeDesc, marshalDesc(desc))

			// 最后升版本号
			exec(conn, "HINCRBY", nodeKey, key_NodeVer, 1)

			// 超时
			exec(conn, "EXPIRE", nodeKey, nodeKeyTimeoutSec)

			ulog.WithField("nodeid", desc.ID).Debugf("register node")

			self.guard.Lock()
			self.nodeByID[desc.ID] = &NodeContext{
				Desc: desc,
				Ver:  1,
			}
			self.dirty = true
			self.guard.Unlock()
		})

	})
}
