package robotutil

import (
	"math/rand"
)

// 整数范围随机，[min, max)
func RandRangeInt32(min, max int32) int32 {
	if min >= max || max == 0 {
		return max
	}
	return min + rand.Int31n(max-min)
}
