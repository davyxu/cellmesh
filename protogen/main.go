package main

import (
	"flag"
	"fmt"
	"github.com/davyxu/cellmesh/protogen/gengo"
	"github.com/davyxu/protoplus/gen"
	"github.com/davyxu/protoplus/model"
	_ "github.com/davyxu/protoplus/msgidutil"
	"github.com/davyxu/protoplus/util"
	"os"
)

var (
	flagPackage = flag.String("package", "", "package name in source files")
	flagGoOut   = flag.String("cmgo_out", "", "cellmesh binding for golang")
)

func main() {

	flag.Parse()

	var err error
	var ctx gen.Context
	ctx.DescriptorSet = new(model.DescriptorSet)
	ctx.PackageName = *flagPackage

	err = util.ParseFileList(ctx.DescriptorSet)

	if err != nil {
		goto OnError
	}

	ctx.OutputFileName = *flagGoOut
	if ctx.OutputFileName != "" {
		err = gengo.GenGo(&ctx)
		if err != nil {
			goto OnError
		}
	}

	return

OnError:
	fmt.Println("Generate error: ", err)
	os.Exit(1)
}
