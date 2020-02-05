package main

import (
	"fmt"
	"github.com/davyxu/protoplus/model"
	"strings"
)

type MsgDir struct {
	From, Mid, To string
	Name          string
}

func (self *MsgDir) HasStar() bool {
	if self.From == "*" {
		return true
	}

	if self.Mid == "*" {
		return true
	}

	if self.To == "*" {
		return true
	}

	return false
}

func (self *MsgDir) Less(Other MsgDir) bool {

	if self.From != Other.From {
		return self.From < Other.From
	}

	if self.Mid != Other.Mid {
		return self.Mid < Other.Mid
	}

	if self.To != Other.To {
		return self.To < Other.To
	}

	return self.Name < Other.Name
}

func ParseMessage(d *model.Descriptor) (rm MsgDir) {

	msgdir := d.TagValueString("MsgDir")
	if msgdir == "" {
		return
	}

	// 上行
	if strings.Contains(msgdir, "->") {
		endPoints := strings.Split(msgdir, "->")
		rm.Name = d.Name

		switch len(endPoints) {
		case 3:
			rm.From = strings.TrimSpace(endPoints[0])

			rm.Mid = strings.TrimSpace(endPoints[1])

			rm.To = strings.TrimSpace(endPoints[2])
			return
		case 2:
			rm.From = strings.TrimSpace(endPoints[0])

			rm.To = strings.TrimSpace(endPoints[1])
			return
		}
	} else if strings.Contains(msgdir, "<-") { // 下行
		endPoints := strings.Split(msgdir, "<-")
		rm.Name = d.Name

		switch len(endPoints) {
		case 3:
			rm.From = strings.TrimSpace(endPoints[2])

			rm.Mid = strings.TrimSpace(endPoints[1])

			rm.To = strings.TrimSpace(endPoints[0])
			return
		case 2:
			rm.From = strings.TrimSpace(endPoints[1])

			rm.To = strings.TrimSpace(endPoints[0])
			return
		}
	} else {
		fmt.Println("unknown msg dir", d.Name, msgdir)
	}

	return
}
