package meshutil

import (
	"sync"
	"time"
)

// numF=半字节
func numFToMask(numF uint) uint64 {
	switch numF {
	case 1:
		return 0xF
	case 2:
		return 0xFF
	case 3:
		return 0xFFF
	case 4:
		return 0xFFFF
	case 5:
		return 0xFFFFF
	case 6:
		return 0xFFFFFF
	case 7:
		return 0xFFFFFFF
	case 8:
		return 0xFFFFFFFF
	case 9:
		return 0xFFFFFFFFF
	case 10:
		return 0xFFFFFFFFFF
	case 11:
		return 0xFFFFFFFFFFF
	case 12:
		return 0xFFFFFFFFFFFF
	case 13:
		return 0xFFFFFFFFFFFFF
	case 14:
		return 0xFFFFFFFFFFFFFF
	case 15:
		return 0xFFFFFFFFFFFFFFF
	case 16:
		return 0xFFFFFFFFFFFFFFFF
	default:
		panic("numF shound in range 1~16")
	}
}

type UUID64Component struct {
	ValueSrc func() uint64 // 数值来源
	NumF     uint          // 占用多少个F 1个F=0.5个字节=4位
}

type UUID64Generator struct {
	seqGen   uint64
	comSet   []*UUID64Component
	genGuard sync.Mutex
}

const (
	MaxNumFInt64 = 16
)

func (self *UUID64Generator) AddComponent(com *UUID64Component) {

	// 检查范围
	numFToMask(com.NumF)

	self.comSet = append(self.comSet, com)

	// 检查总组件超过位数
	if self.UsedNumF() > MaxNumFInt64 {
		panic("total bit over int64(8bit) range")
	}
}

func (self *UUID64Generator) UsedNumF() (ret uint) {

	for _, com := range self.comSet {
		ret += com.NumF
	}

	return
}

func (self *UUID64Generator) LeftNumF() (ret uint) {
	return MaxNumFInt64 - self.UsedNumF()
}

// 序列号组件
func (self *UUID64Generator) AddSeqComponent(numF uint) {
	self.AddComponent(&UUID64Component{
		ValueSrc: func() uint64 {
			self.seqGen++
			return self.seqGen
		},

		NumF: numF,
	})

}

const timeStartPoint = 946656000 // 这里设置参考点为2000/1/1 0:0:0，延迟出现2039年Unix时间戳溢出问题

// 添加时间组件
func (self *UUID64Generator) AddTimeComponent(numF uint) {
	self.AddComponent(&UUID64Component{
		ValueSrc: func() uint64 {
			return uint64(time.Now().Unix() - timeStartPoint)
		},

		NumF: numF,
	})

}

// 按给定的组件规则生成一个UUID
func (self *UUID64Generator) Generate() (ret uint64) {

	self.genGuard.Lock()
	var offset uint
	for _, g := range self.comSet {

		mask := numFToMask(g.NumF)
		part := (g.ValueSrc() & mask) << offset
		ret |= part
		offset += g.NumF * 4
	}

	self.genGuard.Unlock()

	return
}

func NewUUID64Generator() *UUID64Generator {

	return &UUID64Generator{}
}
