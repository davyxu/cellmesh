package db

import (
	"fmt"
	"github.com/davyxu/cellmesh/svc/actor"
	"github.com/gomodule/redigo/redis"
	"reflect"
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

func TestModel(t *testing.T) {

	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ser := NewModelList(c)
	var b actor.Role
	err = ser.Load(&b, int64(2))
	t.Logf("%+v", err)

	var a actor.Role
	a.NickName = "ha"
	a.RoleID = 1
	err = ser.Save(&a, int64(1))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = ser.Load(&b, int64(1))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(a, b) {
		t.FailNow()
	}

	t.Logf("%+v", b)

}

func TestModelBatch(t *testing.T) {

	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	tester := NewModelList(c)

	tester.BatchLoad(new(actor.Role), int64(2))
	tester.Flush()
	for i := 0; i < tester.Count(); i++ {
		_, err := tester.Index(i)

		if err != ErrModelNotExists {
			t.FailNow()
		}
	}

	saver := NewModelList(c)

	for i := 10; i < 12; i++ {

		var a actor.Role
		a.NickName = fmt.Sprintf("a:%d", i)
		a.RoleID = int64(i)
		saver.BatchSave(&a, a.RoleID)
	}

	saver.Flush()

	loader := NewModelList(c)

	for i := 10; i < 12; i++ {

		loader.BatchLoad(new(actor.Role), int64(i))
	}

	loader.Flush()
	for i := 0; i < loader.Count(); i++ {
		_, err := loader.Index(i)

		if err == ErrModelNotExists {
			t.FailNow()
		}

		t.Logf("%+v", loader.MustIndex(i).(*actor.Role))
	}

}
