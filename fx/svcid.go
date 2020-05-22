package fx

import (
	"errors"
	"fmt"
	"strconv"
)

// 全局唯一的svcid 格式:  svcName@ip
// 一台机器开多套时, 需要配合FlagFile指定每组的GroupName, 或者手动参数指定

// 构造指定服务的ID
func MakeSvcID(svcName string) string {
	return fmt.Sprintf("%s#%d@%s", svcName, SvcIndex, SvcGroup)
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
