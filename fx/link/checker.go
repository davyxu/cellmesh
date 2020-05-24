package link

import (
	"fmt"
	"github.com/davyxu/cellmesh/fx"
	"github.com/davyxu/cellmesh/fx/redsd"
	"github.com/davyxu/cellmesh/svc/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/ulog"
	"io"
	"sort"
	"strings"
	"time"
)

func CheckReady() {

	var lastStatus string
	for {

		time.Sleep(time.Second * 3)

		if IsAllReady() {
			ulog.WithColorName("green").Infof("All peers ready!\n%s", localNodeStatus())

			fx.OnLoad.Invoke()

			break
		}

		thisStatus := localNodeStatus()

		if lastStatus != thisStatus {
			ulog.Warnf("\n%s", thisStatus)
			lastStatus = thisStatus
		}

	}
}

func localNodeStatus() string {
	var sb strings.Builder

	nodeListSet := SD.NodeListSet()

	sort.Slice(nodeListSet, func(i, j int) bool {

		a := nodeListSet[i]
		b := nodeListSet[j]

		if a.Kind() != b.Kind() {
			return a.Kind() < b.Kind()
		}

		return a.Name() < b.Name()
	})

	for _, nodeList := range nodeListSet {

		descList := nodeList.DescList()

		sort.Slice(descList, func(i, j int) bool {

			a := descList[i]
			b := descList[j]

			return a.ID < b.ID
		})

		kindStr := proto.NodeKind(nodeList.Kind()).String()

		if len(descList) > 0 {

			for _, desc := range descList {

				printDesc(kindStr, &sb, desc)
			}

		} else {
			fmt.Fprintf(&sb, "%10s| %22s\n", kindStr, nodeList.Name())
		}

	}

	return sb.String()
}

func printDesc(prefix string, w io.Writer, desc *redsd.NodeDesc) {
	fmt.Fprintf(w, "%10s| %22s  %22s %s\n", prefix, desc.ID, desc.Address(), nodeReadyStatus(desc))
}

func isNodeKindShouldReady(kind int) bool {

	switch proto.NodeKind(kind) {
	case proto.NodeKind_Listen, proto.NodeKind_Connect:
		return true
	}

	return false
}

func IsAllReady() bool {

	for _, nodeList := range SD.NodeListSet() {

		if isNodeKindShouldReady(nodeList.Kind()) {
			descList := nodeList.DescList()

			if len(descList) > 0 {

				for _, desc := range descList {
					if !isNodeReady(desc) {
						return false
					}
				}
			} else {
				// 这种类别一个都没有
				return false
			}
		}

	}

	return true
}

func isNodeReady(desc *redsd.NodeDesc) bool {
	if desc.Peer != nil {

		if desc.Peer.(cellnet.PeerReadyChecker).IsReady() {
			return true
		}

	}

	return false
}

func nodeReadyStatus(desc *redsd.NodeDesc) string {

	if isNodeReady(desc) {
		return "[READY]"
	}

	return ""
}

func init() {
	fx.OnCommand.Add(func(args ...interface{}) {

		cmd := args[0]

		switch cmd {
		case "ns", "nodestatus":
			ulog.Debugf("\n%s", localNodeStatus())
		}

	})
}
