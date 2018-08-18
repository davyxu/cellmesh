package gengo

import (
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/protoplus/gen"
	"github.com/davyxu/protoplus/model"
	"sort"
	"text/template"
)

var FuncMap = template.FuncMap{}

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

	FuncMap["HasJsonCodec"] = func(ctx *gen.Context) bool {

		for _, d := range ctx.Structs() {
			if d.TagValueString("Codec") == "json" {
				return true
			}
		}

		return true
	}
}
