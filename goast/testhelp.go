package goast

import (
	"testing"

	"cuelang.org/go/cue/parser"
)

func process(src string, t *testing.T) GoDecls {
	f, err := parser.ParseFile("test", src)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return Build(f)
}
