package memsd

import (
	"github.com/davyxu/cellmesh/discovery"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {

	sd := NewDiscovery(nil)

	sd.SetValue("mykey", 123456)

	var a int
	sd.GetValue("mykey", &a)
	if a != 123456 {
		t.Fatalf("a != 123456")
	}

	sd.DeleteValue("mykey")

	time.Sleep(time.Millisecond * 100)

	err := sd.GetValue("mykey", &a)
	if err == nil {
		t.Fatalf("getvalue == nil")
	}

	if err.Error() != "value not exists" {
		t.Fatalf("value exists")
	}

	sd.Register(&discovery.ServiceDesc{
		Name: "game",
		ID:   "game0@dev",
		Host: "127.0.0.1",
		Port: 9091,
		Tags: []string{"blabla"},
		Meta: map[string]string{"a": "b"},
	})

	time.Sleep(time.Millisecond * 100)

	list := sd.Query("game")
	if len(list) == 0 {
		t.Fatalf("len(list) == 0")
	}

	t.Log(*list[0])

	if list[0].ID != "game0@dev" {
		t.Fatalf("list[0].ID != 'game0@dev'")
	}

	sd.Deregister("game0@dev")

	time.Sleep(time.Millisecond * 100)

	if len(sd.Query("game")) != 0 {
		t.Fatalf("len(sd.Query('game')) != 0")
	}

	time.Sleep(time.Second)
}
