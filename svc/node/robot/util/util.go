package robotutil

import (
	"fmt"
	"github.com/davyxu/cellmesh/svc/proto"
	"github.com/davyxu/cellnet/util"
	"math/rand"
)

func CheckCode(code proto.ResultCode) {
	if code != 0 {
		panic(fmt.Errorf("ErrCode: %s  stack: %s", code.String(), util.StackToString(3)))
	}
}

// 整数范围随机，[min, max)
func RandRangeInt32(min, max int32) int32 {
	if min >= max || max == 0 {
		return max
	}
	return min + rand.Int31n(max-min)
}
