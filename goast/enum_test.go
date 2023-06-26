package goast

import "testing"

func TestEnum1(t *testing.T) {
	src := `#a: int & #b | #c | #d
#b: #a & 1
#c: #a & 2
#d: #a & 3
`

	vals := []*EnumItem{
		{"B", "1"},
		{"C", "2"},
		{"D", "3"},
	}

	gd := process(src, t)

	a := gd.GetEnumeration("A")
	if a == nil {
		t.Error("No enum A")
	}
	for i, item := range a.Items {
		if item.Name != vals[i].Name {
			t.Errorf("item %d Name = %s", i, item.Name)
		}
		if item.Value != vals[i].Value {
			t.Errorf("item %d Value = %s", i, item.Value)
		}
	}
}
