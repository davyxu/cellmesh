package gengo

import (
	"github.com/davyxu/protoplus/gen"
	"github.com/davyxu/protoplus/model"
	"strings"
	"text/template"
)

var FuncMap = template.FuncMap{}

type RPCPair struct {
	REQ *model.Descriptor
	ACK *model.Descriptor
}

func (self *RPCPair) Name() string {
	return strings.TrimSuffix(self.REQ.Name, "REQ")
}

func init() {
	FuncMap["StructCodec"] = func(d *model.Descriptor) string {
		return d.TagValueString("Codec")
	}

	FuncMap["StructService"] = func(d *model.Descriptor) string {
		return d.TagValueString("Service")
	}

	FuncMap["RPCPair"] = func(ctx *gen.Context) (ret []*RPCPair) {

		for _, d := range ctx.Structs() {

			if _, ok := d.TagValueByKey("Service"); ok {

				if strings.HasSuffix(d.Name, "REQ") {

					methodName := strings.TrimSuffix(d.Name, "REQ")

					ack := ctx.ObjectByName(methodName + "ACK")

					if ack != nil {
						ret = append(ret, &RPCPair{REQ: d, ACK: ack})
					}

				}

			}
		}

		return
	}
}
