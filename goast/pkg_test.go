package goast

import "testing"

func TestPkg1(t *testing.T) {
	src := `package pkgtest
	`

	gd := process(src, t)
	if pkg := gd.GetPackage(); pkg != "pkgtest" {
		t.Error(pkg)
	}
}
