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

func init() {
	FuncMap["StructCodec"] = func(d *model.Descriptor) string {
		return d.TagValueString("Codec")
	}

	FuncMap["StructService"] = func(d *model.Descriptor) string {
		return d.TagValueString("Service")
	}

	FuncMap["ProtoImportList"] = func(ctx *gen.Context) (ret []string) {

		linq.From(ctx.Structs()).WhereT(func(d *model.Descriptor) bool {
			return d.TagValueString("Codec") != ""
		}).SelectT(func(d *model.Descriptor) string {
			return d.TagValueString("Codec")
		}).DistinctByT(func(d string) string {
			return d
		}).ToSlice(&ret)

		return
	}

	FuncMap["ServiceGroup"] = func(ctx *gen.Context) (ret []linq.Group) {

		type RecvPair struct {
			Recv string
			d    *model.Descriptor
		}

		var pairs []*RecvPair

		for _, d := range ctx.Structs() {

			recvList := d.TagValueString("Service")

			if recvList == "" {
				continue
			}

			for _, recv := range strings.Split(recvList, "|") {
				pairs = append(pairs, &RecvPair{recv, d})
			}
		}

		linq.From(pairs).GroupByT(func(pair *RecvPair) interface{} {
			return pair.Recv
		}, func(pair *RecvPair) interface{} {
			return pair.d
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
