package main

import (
	_ "github.com/davyxu/cellmesh/demo/agent/backend"
	"github.com/davyxu/cellmesh/demo/agent/frontend"
	"github.com/davyxu/cellmesh/demo/agent/routerule"
	"github.com/davyxu/cellmesh/demo/proto"
	_ "github.com/davyxu/cellmesh/demo/proto" // 进入协议
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellmesh/service/cellsvc"
	"github.com/davyxu/cellmesh/svcfx"
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/golog"
)

var log = golog.New("main")

// TODO 做的和flag一样
func GetKV_String(key string, defaultValue string) (ret string) {

	data, err := discovery.Default.GetValue(key)
	if err != nil {
		ret = defaultValue
		return
	}

	v := string(data)

	if v == "" {
		ret = defaultValue

		// 空值设置
		defer discovery.Default.SetValue(key, []byte(ret))
		return
	}

	ret = v

	return
}

func main() {

	svcfx.Init()

	routerule.Download()

	s := cellsvc.NewService("router")
	s.SetDispatcher(proto.GetDispatcher("router"))
	s.Start()

	frontend.Start(GetKV_String("config/agent/frontend_addr", ":18000"))

	util.WaitExit()

	frontend.Stop()
	s.Stop()
}
