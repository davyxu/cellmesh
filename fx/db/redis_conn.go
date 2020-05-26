package db

import (
	"github.com/davyxu/cellmesh/fx/zonecfg"
	"github.com/davyxu/ulog"
	"github.com/gomodule/redigo/redis"
	"time"
)

type ResultCode = int32

type RedisConn struct {
	pool *redis.Pool
}

func (self *RedisConn) Connect() {
	addr := zonecfg.String("RedisAddress")

	ulog.Debugf("Connecting to redis, %s...", addr)

	self.pool.Dial = func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		// 选择db
		c.Do("SELECT", 0)
		return c, nil
	}
}

func (self *RedisConn) Operate(callback func(conn redis.Conn)) (code ResultCode) {

	conn := self.pool.Get()

	defer func() {

		defer conn.Close()

		switch err := recover().(type) {
		case ResultCode:
			code = err
		case nil:
		default:
			panic(err)
		}

	}()

	callback(conn)

	return
}

func NewRedisConn() *RedisConn {
	return &RedisConn{pool: &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
	}}
}
