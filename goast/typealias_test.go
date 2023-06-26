package goast

import (
	"testing"
)

func TestAlias1(t *testing.T) {
	src := `#a: int`
	ta := process(src, t).GetTypeAlias("A")
	if ta == nil {
		t.Error()
	}
	if ta.name != "A" {
		t.Error(ta.name)
	}
	if ta.GoType != "int" {
		t.Error(ta.GoType)
	}
}

func TestAlias2(t *testing.T) {
	src := `#a: #b
#b : {a : int}`
	ta := process(src, t).GetTypeAlias("A")
	if ta == nil {
		t.Error("ta == nil")
	}
	if ta.GoType != "*B" {
		t.Error(ta.GoType)
	}
}
