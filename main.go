package main

import (
	"flag"
	"fmt"
	"os"

	"cuego/gen"
	"cuego/goast"

	"cuelang.org/go/cue/parser"
)

var (
	cueFile string
	outdir  = flag.String("o", "", "")
)

func main() {
	getParams()

	f, err := parser.ParseFile(cueFile, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	pt := goast.Build(f)
	gen.Gen(cueFile, *outdir, pt)
}

func getParams() {
	flag.Parse()
	if flag.NArg() != 1 {
		fail("exactly one cue file must be specified")
	}
	if *outdir == "" {
		fail("target directory for Go files must be specified")
	}
	cueFile = flag.Arg(0)
}

func fail(fmtstr string, args ...any) {
	fmt.Printf("ERROR: %s\n", fmt.Sprintf(fmtstr, args...))
	fmt.Println(usage)
	os.Exit(1)
}

const usage = `use gostruct -o <dir> <cue file>
where
    <cue file> Mandatory. The input cue file

    -o <dir> Mandatory. The target directory for generated Go files
`
