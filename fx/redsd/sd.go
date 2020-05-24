package redsd

import (
	"github.com/davyxu/cellmesh/fx/db"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/ulog"
	"github.com/gomodule/redigo/redis"
	"sync"

	"time"
)

type NodeContext struct {
	Desc *NodeDesc
	Ver  int64
}

type RedisDiscovery struct {
	pool *redis.Pool

	nodeByID sync.Map

	looper *time.Ticker

	nodeListByName      map[string]*NodeList
	nodeListByNameGuard sync.RWMutex
}

func (self *RedisDiscovery) operate(callback func(conn redis.Conn)) {

	conn := self.pool.Get()

	defer conn.Close()

	callback(conn)
}

func exec(conn redis.Conn, commandName string, args ...interface{}) error {

	_, err := conn.Do(commandName, args...)
	if err != nil {
		ulog.Errorf("redis do failed, %s , stack: %s", err, util.StackToString(5))
	}

	return err
}

func (self *RedisDiscovery) NodeListByName(name string) *NodeList {
	self.nodeListByNameGuard.RLock()
	defer self.nodeListByNameGuard.RUnlock()
	c, _ := self.nodeListByName[name]
	return c
}

func (self *RedisDiscovery) NodeListSet() (ret []*NodeList) {
	self.nodeListByNameGuard.RLock()
	defer self.nodeListByNameGuard.RUnlock()

	for _, nodeList := range self.nodeListByName {
		ret = append(ret, nodeList)
	}

	return
}

func (self *RedisDiscovery) NewNodeList(name string, kind int) *NodeList {

	c := newNodeList(name, self)
	c.kind = kind

	self.nodeListByNameGuard.Lock()

	if _, ok := self.nodeListByName[name]; ok {
		panic("duplicate node name: " + name)
	}

	self.nodeListByName[name] = c
	self.nodeListByNameGuard.Unlock()

	return c
}

func safeOpNode(conn redis.Conn, nodeid string, callback func()) {
	lock := db.NewRedisLock(conn, key_NodeLock+nodeid)

	if err := lock.Lock(); err != nil {
		ulog.WithField("nodeid", nodeid).Errorf("node op failed, %s, stack: %s", err, util.StackToString(5))
	} else {
		callback()
	}

	lock.Unlock()
}

func (self *RedisDiscovery) checkRegNode() {

	self.nodeListByNameGuard.RLock()

	var regNodeList []*NodeDesc
	for _, nodeList := range self.nodeListByName {
		descList := nodeList.DescList()
		if nodeList.heartBeat && len(descList) > 0 {
			regNodeList = append(regNodeList, descList[0])
		}
	}

	self.nodeListByNameGuard.RUnlock()

	self.operate(func(conn redis.Conn) {

		for _, desc := range regNodeList {
			nodeKey := key_Node + desc.ID

			safeOpNode(conn, desc.ID, func() {

				//ulog.WithField("nodeid", desc.ID).Debugf("heartbeat node")

				// 心跳
				exec(conn, "EXPIRE", nodeKey, nodeKeyTimeoutSec)
			})
		}

	})
}

func (self *RedisDiscovery) Stop() {
	ulog.Debugf("redis discovery stop")
	self.looper.Stop()
	self.pool.Close()
}

func (self *RedisDiscovery) Start(addr string) {

	self.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", 0)
			return c, nil
		},
	}

	self.looper = time.NewTicker(nodeKeyCheckIntervalSec * time.Second)

	go func() {
		for {
			<-self.looper.C
			self.checkRegNode()
		}
	}()

	ulog.Debugf("redis discovery ready")

}

func NewRedisDiscovery() *RedisDiscovery {
	return &RedisDiscovery{
		nodeListByName: make(map[string]*NodeList),
	}
}
