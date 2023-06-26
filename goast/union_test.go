package goast

import "testing"

func TestUnion1(t *testing.T) {
	src := `#a : int8 | #b | #c
	#b : int8
	#c : {
		a : bool
	}
`
	fieldSpec := [][]string{
		{"Int8", "int8"},
		{"B", "B"},
		{"C", "*C"},
	}
	gd := process(src, t)
	s := gd.GetUnionMessage("A")

	for i, fs := range fieldSpec {
		if s.Fields[i].Name != fs[0] {
			t.Errorf("Field %d: expected field name %s, got %s", i, fs[0], s.Fields[i].Name)
		}
		if s.Fields[i].GoType != fs[1] {
			t.Errorf("Field %d: expected field type %s, got %s", i, fs[1], s.Fields[i].GoType)
		}
	}
}
