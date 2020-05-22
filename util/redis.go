package meshutil

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	ErrLockFailed = errors.New("redis lock failed")
)

// Redis分布式锁
type RedisLock struct {
	conn         redis.Conn
	key          string
	KeyExpireSec int
	LockTimeout  time.Duration
	FailWait     time.Duration
}

func (self *RedisLock) Lock() error {

	beginLock := time.Now()

	for {

		raw, err := self.conn.Do("SET", self.key, 1, "EX", self.KeyExpireSec, "NX")

		if err != nil {
			return err
		}

		switch result := raw.(type) {
		case string:
			if result == "OK" {
				return nil
			}
		}

		if time.Since(beginLock) > self.LockTimeout {
			break
		} else {
			time.Sleep(self.FailWait)
		}
	}

	return ErrLockFailed
}

func (self *RedisLock) Unlock() error {
	_, err := self.conn.Do("DEL", self.key)
	return err
}

func NewRedisLock(conn redis.Conn, key string) *RedisLock {
	return &RedisLock{
		conn:         conn,
		key:          key,
		LockTimeout:  time.Second * 2,
		KeyExpireSec: 3,
		FailWait:     time.Millisecond * 300,
	}
}
