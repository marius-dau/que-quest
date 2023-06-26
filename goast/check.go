package goast

func (bld *builder) checkPackage() {
	if bld.goDecls.GetPackage() == "" {
		fail(bld.f.Pos(), "go package declaration is required")
	}
}
