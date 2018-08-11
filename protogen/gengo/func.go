package gengo

import (
	"github.com/davyxu/protoplus/model"
	"text/template"
)

var FuncMap = template.FuncMap{}

func init() {
	FuncMap["StructCodec"] = func(d *model.Descriptor) string {
		return d.TagValueString("Codec")
	}
}
