package link

import (
	"github.com/davyxu/cellmesh/proto"
	"time"
)

// 监控某个服务的状态
func MonitorService(nodeName string, duration time.Duration) {

	nodeList := SD.NewNodeList(nodeName, int(proto.NodeKind_Monitor))
	nodeList.Monitor(duration)
}
