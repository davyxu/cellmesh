package gengo

import (
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/protoplus/gen"
	"github.com/davyxu/protoplus/model"
	"sort"
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

	FuncMap["ServiceGroup"] = func(ctx *gen.Context) (ret []linq.Group) {

		linq.From(ctx.Structs()).WhereT(func(d *model.Descriptor) bool {
			return d.TagValueString("Service") != ""
		}).GroupByT(func(d *model.Descriptor) interface{} {
			return d.TagValueString("Service")
		}, func(d *model.Descriptor) interface{} {
			return d
		}).SortT(func(a, b linq.Group) bool {

			asvc := a.Key.(string)
			bsvc := b.Key.(string)

			return asvc < bsvc
		}).ToSlice(&ret)

		for _, g := range ret {
			sort.Slice(g.Group, func(i, j int) bool {
				a := g.Group[i].(*model.Descriptor)
				b := g.Group[j].(*model.Descriptor)

				return a.Name < b.Name
			})
		}

		return
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

	FuncMap["HasJsonCodec"] = func(ctx *gen.Context) bool {

		for _, d := range ctx.Structs() {
			if d.TagValueString("Codec") == "json" {
				return true
			}
		}

		return true
	}
}
