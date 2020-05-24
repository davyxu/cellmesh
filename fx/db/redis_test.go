package db

import (
	"github.com/gomodule/redigo/redis"
	"testing"
)

func TestLock(t *testing.T) {

	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	l := NewRedisLock(c, "ttt")

	if err := l.Lock(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := l.Lock(); err != ErrLockFailed {
		t.Error(err)
		t.FailNow()
	}

	l.Unlock()

	if err := l.Lock(); err != nil {
		t.Error(err)
		t.FailNow()
	}

}
