package hubstatus

import (
	"github.com/davyxu/cellmesh/demo/svc/hub/model"
	"github.com/davyxu/cellmesh/service"
	"math/rand"
	"sort"
)

func SelectServiceByLowUserCount(svcName, svcGroup string, mustConnected bool) (finalSvcID string) {

	var statusList []*model.Status
	model.VisitStatus(func(status *model.Status) bool {

		name, _, group, err := service.ParseSvcID(status.SvcID)
		if err != nil {
			return true
		}

		if name != svcName {
			return true
		}

		if svcGroup == "" || svcGroup == group {

			if !mustConnected || service.GetRemoteService(status.SvcID) != nil {
				statusList = append(statusList, status)
			}

		}

		return true
	})

	total := len(statusList)

	switch total {
	case 0:
		return ""
	case 1:
		return statusList[0].SvcID
	default:

		sort.Slice(statusList, func(i, j int) bool {
			a := statusList[i]
			b := statusList[j]
			return a.UserCount < b.UserCount
		})

		lowRange := MaxInt32(int32(total/3), int32(total))
		lowList := statusList[:lowRange]

		final := lowList[rand.Int31n(lowRange)]
		finalSvcID = final.SvcID

	}

	return

}

func MaxInt32(a, b int32) int32 {
	if a > b {
		return a
	}

	return b
}
