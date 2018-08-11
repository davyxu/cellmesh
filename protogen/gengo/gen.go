package gengo

import (
	"fmt"
	"github.com/davyxu/protoplus/codegen"
	"github.com/davyxu/protoplus/gen"
)

func GenGo(ctx *gen.Context) error {

	gen := codegen.NewCodeGen("cmgo").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(FuncMap).
		ParseTemplate(goCodeTemplate, ctx).
		FormatGoCode()

	if gen.Error() != nil {
		fmt.Println(string(gen.Data()))
		return gen.Error()
	}

	return gen.WriteOutputFile(ctx.OutputFileName).Error()
}
