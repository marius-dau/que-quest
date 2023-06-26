package gen

import (
	"bytes"
	"cuego/goast"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/goccmack/goutil/ioutil"
)

func Gen(infile, outdir string, gd goast.GoDecls) {
	w := new(bytes.Buffer)
	genPackage(w, gd)
	genGoDecls(w, gd)

	gocode := w.Bytes()
	if err := ioutil.WriteFile(getOutFile(infile, outdir)); err != nil {
		panic(err)
	}
}

func genGoDecls(w *bytes.Buffer, gd goast.GoDecls) {
	for _, decl := range gd {
		switch d := decl.(type) {
		case *goast.Enumeration:
			genEnumeration(w, d)
		case *goast.StructMessage:
			genStructMessage(w, d)
		case *goast.Package:
			// ignore
		case *goast.TypeAlias:
			genTypeAlias(w, d)
		case *goast.TypedConst:
			genTypedConst(w, d)
		case *goast.UnionMessage:
			genUnionMessage(w, d)
		default:
			panic("impossible")
		}
	}
}

func genEnumeration(w *bytes.Buffer, d *goast.Enumeration) {
	fmt.Fprintf(w, "type %s %s\n\n", d.Name(), d.GoType)
	fmt.Fprint(w, "const (\n")
	for _, item := range d.Items {
		fmt.Fprintf(w, "    %s %s = %s\n", item.Name, d.Name(), item.Value)
	}
	fmt.Fprint(w, ")\n\n")
}

func genStructMessage(w *bytes.Buffer, d *goast.StructMessage) {
	fmt.Fprintf(w, "type %s struct {\n", d.Name())
	for _, fld := range d.Fields {
		genMessageField(w, fld)
	}
	fmt.Fprint(w, "}\n\n")
}

func genPackage(w *bytes.Buffer, gd goast.GoDecls) {
	fmt.Fprintf(w, "package %s\n\n", gd.GetPackage())
}

func genTypeAlias(w *bytes.Buffer, d *goast.TypeAlias) {
	fmt.Fprintf(w, "implement %T\n\n", d)
}

func genTypedConst(w *bytes.Buffer, d *goast.TypedConst) {
	fmt.Fprintf(w, "const %s %s = %s\n\n", d.Name(), d.GoType, d.Value)
}

func genUnionMessage(w *bytes.Buffer, d *goast.UnionMessage) {
	fmt.Fprintf(w, "type %s struct {\n", d.Name())
	fmt.Fprintf(w, "    MsgType %sType\n", d.Name())
	for _, fld := range d.Fields {
		genMessageField(w, fld)
	}
	fmt.Fprint(w, "}\n\n")
	genUnionMessageTypes(w, d)
}

func genUnionMessageTypes(w *bytes.Buffer, d *goast.UnionMessage) {
	fmt.Fprintf(w, "type %sType int16\n\n", d.Name())
	fmt.Fprintf(w, "const(\n")
	for i, fld := range d.Fields {
		fmt.Fprintf(w, "    %sType %sType = %d\n", fld.Name, d.Name(), i+1)
	}
	fmt.Fprint(w, ")\n")
}

func genMessageField(w *bytes.Buffer, fld *goast.MessageField) {
	fmt.Fprintf(w, "    %s %s\n", fld.Name, fld.GoType)
}

func getOutFile(infile, outdir string) string {
	_, infilename := filepath.Split(infile)
	fnames := strings.Split(infilename, ".")
	outFileName := strings.Join(fnames[:len(fnames)-1], ".")
	return filepath.Join(outdir, outFileName+".go")
}
