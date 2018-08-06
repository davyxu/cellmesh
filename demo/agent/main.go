package main

import (
	"github.com/davyxu/cellmesh/demo/agent/router"
	"github.com/davyxu/cellmesh/util"
)

func main() {
	router.Start()

	util.WaitExit()

	router.Stop()
}
