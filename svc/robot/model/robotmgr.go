package model

import (
	"fmt"
	"github.com/davyxu/cellnet/util"
	"net"
	"sync"
)

var (
	robotByID sync.Map

	allStates      []string
	allStatesGuard sync.RWMutex
)

func AddRobot(robot *Robot) {
	robotByID.Store(robot.ID, robot)
}

func VisitRobot(callback func(robot *Robot) bool) {

	robotByID.Range(func(key, value interface{}) bool {
		return callback(value.(*Robot))
	})
}

func AddState(state string) {

	allStatesGuard.Lock()
	defer allStatesGuard.Unlock()

	for _, s := range allStates {
		if s == state {
			return
		}
	}
	allStates = append(allStates, state)
}

func VisitAllState(callback func(string)) {

	allStatesGuard.RLock()
	for _, s := range allStates {
		callback(s)
	}
	allStatesGuard.RUnlock()
}

func GenBaseID() string {

	localIP := util.GetLocalIP()
	ipBytes := net.ParseIP(localIP)

	num1 := int64(ipBytes[len(ipBytes)-1])
	num2 := int64(ipBytes[len(ipBytes)-2])

	return fmt.Sprintf("%d%d", num2, num1)
}
