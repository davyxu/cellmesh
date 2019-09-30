package cellmesh

import (
	"errors"
	"fmt"
	"strconv"
)

// 全局唯一的svcid 格式:  svcName#svcIndex@svcGroup

// 构造服务ID
func MakeSvcID(svcName string, svcIndex int, svcGroup string) string {
	return fmt.Sprintf("%s#%d@%s", svcName, svcIndex, svcGroup)
}

// 构造指定服务的ID
func MakeLocalSvcID(svcName string) string {
	return MakeSvcID(svcName, *flagSvcIndex, *flagSvcGroup)
}

// 获得本进程的服务id
func GetLocalSvcID() string {
	return MakeLocalSvcID(GetProcName())
}

func ParseSvcID(svcid string) (svcName string, svcIndex int, svcGroup string, err error) {

	var sharpPos, atPos = -1, -1

	for pos, c := range svcid {
		switch c {
		case '#':
			sharpPos = pos
			svcName = svcid[:sharpPos]
		case '@':
			atPos = pos
			svcGroup = svcid[atPos+1:]

			if sharpPos == -1 {
				break
			}

			var n int64
			n, err = strconv.ParseInt(svcid[sharpPos+1:atPos], 10, 32)
			if err != nil {
				break
			}
			svcIndex = int(n)
		}
	}

	if sharpPos == -1 {
		err = errors.New("missing '#' in svcid")
	}

	if atPos == -1 {
		err = errors.New("missing '@' in svcid")
	}

	return
}
