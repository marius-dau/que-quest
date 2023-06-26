package goast

import (
	"testing"
)

func Test1(t *testing.T) {
	src := `#a : int8 & 5`
	goDecls := process(src, t)
	c := goDecls.GetConst("A")
	if c == nil {
		t.Error()
	}
	if c.GoType != "int8" {
		t.Error()
	}
	if c.Value != "5" {
		t.Error()
	}
}

func Test2(t *testing.T) {
	src := `#a : #b & 5
#b: int8`

	goDecls := process(src, t)
	c := goDecls.GetConst("A")
	if c == nil {
		t.Error()
	}
	if c.GoType != "B" {
		t.Error(c.GoType)
	}
	if c.Value != "5" {
		t.Error("Value:", c.Value)
	}
	cb := goDecls.GetTypeAlias("B")
	if cb == nil {
		t.Error("No alias B declared")
	}
}
