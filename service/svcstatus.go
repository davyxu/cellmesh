package service

import (
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/timer"
	"reflect"
	"strings"
	"svc/table"
	"time"
)

func MakeConfigKey() string {
	return fmt.Sprintf("status/%s(%s)", GetProcName(), GetNode())
}

func ParseConfigKey(key string) (svcid string) {

	pathIndex := strings.Index(key, "/")
	if pathIndex == -1 {
		return
	}

	svcid = key[pathIndex+1:]

	sqIndex := strings.Index(svcid, "(")
	if sqIndex == -1 {
		return
	}

	var id svctable.SvcID
	id.Name = svcid[:sqIndex]

	id.Node = svcid[sqIndex+1 : len(svcid)-1]

	return id.String()
}

// 定时汇报状况
func StartSendStatus(q cellnet.EventQueue, interval time.Duration, statusCallback func() interface{}) {

	timer.NewLoop(q, interval, func(loop *timer.Loop) {

		if discovery.Default != nil {
			discovery.Default.SetValue(MakeConfigKey(), statusCallback())
		}

	}, nil).Notify().Start()
}

func QueryServiceStatus(svcName string, statusType reflect.Type, callback func(svcid string, status interface{}) bool) error {
	valueList, err := discovery.Default.GetRawValueList("status/" + svcName)
	if err != nil {
		return err
	}

	for _, value := range valueList {

		svcid := ParseConfigKey(value.Key)

		dataPtr := reflect.New(statusType).Interface()

		if err := discovery.BytesToAny(value.Value, dataPtr); err != nil {
			return err
		}

		if !callback(svcid, dataPtr) {
			return nil
		}
	}

	return nil
}
