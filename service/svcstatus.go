package service

import (
	"fmt"
	"github.com/davyxu/cellmesh/discovery"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/timer"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func MakeConfigKey() string {
	return fmt.Sprintf("status/%s.%d.%s", GetProcName(), GetSvcIndex(), GetSvcGroup())
}

func ParseConfigKey(key string) (svcid string) {

	pathIndex := strings.Index(key, "/")
	if pathIndex == -1 {
		return
	}

	svcid = key[pathIndex+1:]

	triples := strings.Split(svcid, ".")
	if len(triples) != 3 {
		return
	}

	svcIndex, err := strconv.ParseInt(triples[1], 10, 32)
	if err != nil {
		return
	}

	return MakeSvcID(triples[0], int(svcIndex), triples[2])
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
