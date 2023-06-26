package gen

import "testing"

func TestGen1(t *testing.T) {
	out := getOutFile("a/b/c.cue", "d")
	if out != "d/c.go" {
		t.Error(out)
	}
}

func TestGen2(t *testing.T) {
	out := getOutFile("a/b/c.d.e.cue", "f")
	if out != "f/c.d.e.go" {
		t.Error(out)
	}
}
