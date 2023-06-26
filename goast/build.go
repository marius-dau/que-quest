package goast

import (
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/token"

	"github.com/iancoleman/strcase"
)

type builder struct {
	f       *ast.File
	goDecls GoDecls
}

func Build(f *ast.File) GoDecls {
	bld := &builder{
		f: f,
	}
	bld.getGoDecls(f.Decls)
	bld.setPointerTypes()
	bld.checkPackage()
	return bld.goDecls
}

func (bld *builder) addEnumItem(f *ast.Field) {
	itemName := getGoFieldName(f)
	binExpr := f.Value.(*ast.BinaryExpr)
	enumID := binExpr.X.(*ast.Ident)
	val := binExpr.Y.(*ast.BasicLit).Value
	item := bld.goDecls.GetEnumeration(goFieldName(enumID.Name)).GetItem(itemName)
	if item == nil {
		fail(f.TokenPos, "enumeration %s had no item %s", enumID.Name, itemName)
	}
	item.Value = val
}

func (bld *builder) getGoDecl(decl ast.Decl) GoDecl {
	switch d := decl.(type) {
	case *ast.Package:
		return bld.getPackage(d)
	case *ast.Field:
		// fmt.Println(fmt.Sprintf("%T %s %T", d, d.Label, d.Value))
		switch {
		case isTypeAlias(d):
			return bld.getTypeAlias(d)
		case bld.isEnumItem(d):
			bld.addEnumItem(d)
			return nil
		case isTypedConst(d):
			return bld.getTypedConst(d)
		case isUnionMessage(d):
			return bld.getUnionMessage(d)
		case isStructMessage(d):
			return bld.getStructMessage(d)
		case isEnumeration(d):
			return bld.getEnumeration(d)
		}
		fail(d.Pos(), "unsupported cue statement")
	}
	panic("impossible")
}

func (bld *builder) getGoDecls(decls []ast.Decl) {
	for _, d := range decls {
		if decl := bld.getGoDecl(d); decl != nil {
			bld.goDecls = bld.goDecls.Add(decl)
		}
	}
	return
}

func (bld *builder) getEnumeration(f *ast.Field) *Enumeration {
	ids, _ := getBinaryExprIdents(f.Value.(*ast.BinaryExpr))
	enum := &Enumeration{
		name:   getGoFieldName(f),
		GoType: bld.getGoType(ids[0]),
	}
	for _, item := range ids[1:] {
		itemName := goFieldName(item.Name)
		if nil != bld.goDecls.GetDecl(itemName) {
			fail(item.Pos(), "item name %s is already declared", item.Name)
		}
		if nil != enum.GetItem(itemName) {
			fail(item.Pos(), "duplicate enum item %s", item.Name)
		}
		enum.Items = append(enum.Items,
			&EnumItem{
				Name: itemName,
				// Value is declared separately
			})
	}
	return enum
}

func (bld *builder) getPackage(p *ast.Package) *Package {
	return &Package{
		name: p.Name.Name,
	}
}

func (bld *builder) getGoType(id *ast.Ident) string {
	if got, exist := cueTypeToGo[id.Name]; exist {
		return got
	}
	f := getField(bld.f, id.Name)
	if f == nil {
		fail(id.NamePos, "not builtin or declared type: %s", id.Name)
	}
	typeName := goFieldName(id.Name)
	if isStructMessage(f) {
		return "*" + typeName
	}
	return typeName
}

func (bld *builder) getMessageField(decl ast.Decl) *MessageField {
	fld, ok := decl.(*ast.Field)
	if !ok {
		fail(decl.Pos(), "field expected")
	}
	return &MessageField{
		Name:   goFieldName(getGoFieldName(fld)),
		GoType: bld.getMessageFieldType(fld),
	}
	panic("implement")
}

func (bld *builder) getMessageFieldType(fld *ast.Field) string {
	var typeName *ast.Ident
	switch fldVal := fld.Value.(type) {
	case *ast.Ident:
		typeName = fldVal
	case *ast.BinaryExpr:
		typeNames, err := getBinaryExprIdents(fldVal)
		if err != nil {
			fail(fld.TokenPos, "first element of bin expr is not an ident")
		}
		typeName = typeNames[0]
	default:
		fail(fld.Value.Pos(), "%T", fld.Value)
	}
	return goType(typeName.Name)
}

func (bld *builder) getStructMessage(f *ast.Field) *StructMessage {
	sm := &StructMessage{
		name: getGoFieldName(f),
	}
	sl, ok := f.Value.(*ast.StructLit)
	if !ok {
		fail(f.TokenPos, "StructLit expected")
	}
	for _, decl := range sl.Elts {
		var err error
		sm.Fields, err = sm.Fields.Add(bld.getMessageField(decl))
		if err != nil {
			fail(decl.Pos(), "%s", err)
		}
	}
	return sm
}

func (bld *builder) getTypeAlias(f *ast.Field) *TypeAlias {
	return &TypeAlias{
		name:   getGoFieldName(f),
		GoType: bld.getGoType(f.Value.(*ast.Ident)),
	}
	panic("implement")
}

func (bld *builder) getTypedConst(f *ast.Field) *TypedConst {
	binExpr := f.Value.(*ast.BinaryExpr)
	tc := &TypedConst{
		name:   getGoFieldName(f),
		GoType: bld.getGoType(binExpr.X.(*ast.Ident)),
		Value:  binExpr.Y.(*ast.BasicLit).Value,
	}
	if strings.HasPrefix(tc.GoType, "*") {
		fail(f.TokenPos, "const cannot have a struct type")
	}
	return tc
}

func (bld *builder) getUnionMessage(f *ast.Field) *UnionMessage {
	msg := &UnionMessage{
		name: getGoFieldName(f),
	}
	idents, _ := getBinaryExprIdents(f.Value.(*ast.BinaryExpr))
	for _, id := range idents {
		var err error
		msg.Fields, err = msg.Fields.Add(bld.getUnionMessageField(id))
		if err != nil {
			fail(id.NamePos, err.Error())
		}
	}
	return msg
}

func (bld *builder) getUnionMessageField(id *ast.Ident) *MessageField {
	return &MessageField{
		Name:   goFieldName(id.Name),
		GoType: goType(id.Name),
	}
}

func (bld *builder) setPointerTypes() {
	for _, decl := range bld.goDecls {
		switch d := decl.(type) {
		case *Enumeration:
			// Do nothing
		case *Package:
			// Do nothing
		case *StructMessage:
			bld.setStructMessagePointerTypes(d)
		case *TypeAlias:
			bld.setTypeAliasPointerTypes(d)
		case *TypedConst:
			// Do nothing
		case *UnionMessage:
			bld.setUnionMessagePointerTypes(d)
		default:
			panic("impossible")
		}
	}
}

func (bld *builder) isPointerType(goType string) bool {
	if bld.goDecls.GetStructMessage(goType) != nil {
		return true
	}
	return false
}

func (bld *builder) setTypeAliasPointerTypes(s *TypeAlias) {
	if bld.isPointerType(s.GoType) {
		s.GoType = "*" + s.GoType
	}
}

func (bld *builder) setStructMessagePointerTypes(s *StructMessage) {
	for _, fld := range s.Fields {
		if bld.isPointerType(fld.GoType) {
			fld.GoType = "*" + fld.GoType
		}
	}
}

func (bld *builder) setUnionMessagePointerTypes(s *UnionMessage) {
	for _, fld := range s.Fields {
		if bld.isPointerType(fld.GoType) {
			fld.GoType = "*" + fld.GoType
		}
	}
}

//***** utils *****

func getBinaryExprIdents(be *ast.BinaryExpr) (idents []*ast.Ident, err error) {
	if be1, isBE := be.X.(*ast.BinaryExpr); isBE {
		idents, err = getBinaryExprIdents(be1)
		if err != nil {
			return nil, err
		}
	} else {
		id, ok := be.X.(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("not ident")
		}
		idents = []*ast.Ident{id}
	}
	if be2, isBE := be.Y.(*ast.BinaryExpr); isBE {
		ids, err := getBinaryExprIdents(be2)
		if err != nil {
			return nil, err
		}
		idents = append(idents, ids...)
	} else {
		id, ok := be.Y.(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("not ident")
		}
		idents = append(idents, id)
	}
	return
}

func getBinaryExprOperators(be *ast.BinaryExpr) (ops []token.Token) {
	if be1, isBE := be.X.(*ast.BinaryExpr); isBE {
		ops = getBinaryExprOperators(be1)
	}
	ops = append(ops, be.Op)
	if be2, isBE := be.Y.(*ast.BinaryExpr); isBE {
		ops = append(ops, getBinaryExprOperators(be2)...)
	}
	return
}

func getField(f *ast.File, declName string) *ast.Field {
	for _, d := range f.Decls {
		if fld, ok := d.(*ast.Field); ok {
			if getFieldName(fld) == declName {
				return fld
			}
		}
	}
	return nil
}

func getFieldName(fld *ast.Field) string {
	return fmt.Sprintf("%s", fld.Label)
}

func goFieldName(id string) string {
	fname := strings.TrimPrefix(id, "#")
	return strcase.ToCamel(fname)
}

func getGoFieldName(fld *ast.Field) string {
	return goFieldName(getFieldName(fld))
}

var cueTypeToGo = map[string]string{
	"bool":    "bool",
	"int":     "int",
	"int8":    "int8",
	"int16":   "int16",
	"int32":   "int32",
	"int64":   "int64",
	"uint8":   "uint8",
	"uint16":  "uint16",
	"uint32":  "uint32",
	"uint64":  "uint64",
	"float32": "float32",
	"float64": "float64",
	"string":  "string",
}

var isGoBaseType = func() map[string]bool {
	btMap := map[string]bool{}
	for _, goType := range cueTypeToGo {
		btMap[goType] = true
	}
	return btMap
}()

func goType(id string) string {
	id = strings.TrimPrefix(id, "#")
	if got, exists := cueTypeToGo[id]; exists {
		return got
	}
	return goFieldName(id)
}

func isTypedConst(f *ast.Field) bool {
	binExpr, ok := f.Value.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	_, ok = binExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	if binExpr.Op != token.AND {
		return false
	}
	_, ok = binExpr.Y.(*ast.BasicLit)
	return ok
}

func isIdent(expr ast.Expr) bool {
	_, isIdent := expr.(*ast.Ident)
	return isIdent
}

func (bld *builder) isEnumItem(f *ast.Field) bool {
	binExpr, ok := f.Value.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	if binExpr.Op != token.AND {
		return false
	}
	enumID, ok := binExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	if nil == bld.goDecls.GetEnumeration(goFieldName(enumID.Name)) {
		return false
	}
	_, ok = binExpr.Y.(*ast.BasicLit)
	return ok
}

func isEnumType(id *ast.Ident) bool {
	return id.Name == "int" || id.Name == "string"
}

func isEnumeration(f *ast.Field) bool {
	binExp, ok := f.Value.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	ops := getBinaryExprOperators(binExp)
	if ops[0] != token.AND {
		return false
	}
	idents, err := getBinaryExprIdents(binExp)
	if err != nil {
		return false
	}
	if !isEnumType(idents[0]) {
		return false
	}
	for _, op := range ops[1:] {
		if op != token.OR {
			return false
		}
	}
	return true
}

func isStructMessage(f *ast.Field) bool {
	_, ok := f.Value.(*ast.StructLit)
	return ok
}

func isTypeAlias(f *ast.Field) bool {
	_, isTypeAlias := f.Value.(*ast.Ident)
	return isTypeAlias
}

func isUnionMessage(f *ast.Field) bool {
	binExp, ok := f.Value.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	for _, op := range getBinaryExprOperators(binExp) {
		if op != token.OR {
			return false
		}
	}
	_, err := getBinaryExprIdents(binExp)
	if err != nil {
		return false
	}
	return true
}

func fail(pos token.Pos, fmtstr string, args ...any) {
	msg := fmt.Sprintf(fmtstr, args...)
	fmt.Fprintf(os.Stderr, "error %s: %s\n", pos, msg)
	os.Exit(1)
}
