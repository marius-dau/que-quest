package goast

import "testing"

func TestStruct1(t *testing.T) {
	src := `#a : {
		a : int8
		b : #b
		c : #c
	}
	
	#b : int8
	#c : {
		a : bool
	}
`
	fieldSpec := [][]string{
		{"A", "int8"},
		{"B", "B"},
		{"C", "*C"},
	}
	gd := process(src, t)
	s := gd.GetStructMessage("A")

	for i, fs := range fieldSpec {
		if s.Fields[i].Name != fs[0] {
			t.Errorf("Field %d: expected field name %s, got %s", i, fs[0], s.Fields[i].Name)
		}
		if s.Fields[i].GoType != fs[1] {
			t.Errorf("Field %d: expected field type %s, got %s", i, fs[1], s.Fields[i].GoType)
		}
	}

}

func TestStruct2(t *testing.T) {
	src := `#a : {
		b : int8 & #c
	}
	#c : int8 & 3`
	gd := process(src, t)
	sm := gd.GetStructMessage("A")
	if sm.Fields[0].GoType != "int8" {
		t.Error(sm.Fields[0].GoType)
	}
}
