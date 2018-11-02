package model

import (
	"time"
)

type Status struct {
	UserCount int32

	SvcID string

	LastUpdate time.Time
}

var (
	statusBySvcID = map[string]*Status{}
)

func UpdateStatus(nowStatus *Status) *Status {

	status, _ := statusBySvcID[nowStatus.SvcID]
	if status == nil {
		status = nowStatus
		statusBySvcID[nowStatus.SvcID] = status
	}

	status.UserCount = nowStatus.UserCount
	status.LastUpdate = time.Now()

	return status
}

func RemoveStatus(svcid string) {
	delete(statusBySvcID, svcid)
}

func VisitStatus(callback func(status *Status) bool) {

	for _, s := range statusBySvcID {
		if !callback(s) {
			break
		}
	}
}
