package goast

import (
	"fmt"
)

type GoDecl interface {
	isGoDecl()
	Name() string
}

func (*Enumeration) isGoDecl()   {}
func (*StructMessage) isGoDecl() {}
func (*Package) isGoDecl()       {}
func (*TypeAlias) isGoDecl()     {}
func (*TypedConst) isGoDecl()    {}
func (*UnionMessage) isGoDecl()  {}

type GoDecls []GoDecl

type Enumeration struct {
	name   string
	GoType string
	Items  []*EnumItem
}

type EnumItem struct {
	Name  string
	Value string
}

type MessageField struct {
	Name       string
	GoType     string
	IsPointer  bool
	Constraint string
}

type MessageFields []*MessageField

type Package struct {
	name string
}

type StructMessage struct {
	name   string
	Fields MessageFields
}

type TypeAlias struct {
	name   string
	GoType string

	// typeVal *ast.Ident
}

type TypedConst struct {
	name   string
	GoType string
	Value  string
}

type UnionMessage struct {
	name   string
	Fields MessageFields
}

func (enum *Enumeration) GetItem(name string) *EnumItem {
	for _, item := range enum.Items {
		if item.Name == name {
			return item
		}
	}
	return nil
}

func (d *Enumeration) Name() string {
	return d.name
}

func (d *Package) Name() string {
	return d.name
}

func (d *StructMessage) Name() string {
	return d.name
}

func (d *TypedConst) Name() string {
	return d.name
}

func (d *TypeAlias) Name() string {
	return d.name
}

func (d *UnionMessage) Name() string {
	return d.name
}

//***** GoDecls methods *****

func (fs MessageFields) Add(f *MessageField) (MessageFields, error) {
	if fs.Contain(f.Name) {
		return nil, fmt.Errorf("duplicate field name: %s", f.Name)
	}
	return append(fs, f), nil
}

func (fs MessageFields) Contain(name string) bool {
	for _, f := range fs {
		if f.Name == name {
			return true
		}
	}
	return false
}

func (g GoDecls) Contains(name string) bool {
	for _, d := range g {
		if d.Name() == name {
			return true
		}
	}
	return false
}

func (g GoDecls) Add(d GoDecl) GoDecls {
	return append(g, d)
}

func (g GoDecls) GetConst(name string) *TypedConst {
	for _, d := range g {
		if c, ok := d.(*TypedConst); ok {
			if c.name == name {
				return c
			}
		}
	}
	return nil
}

func (g GoDecls) GetDecl(name string) GoDecl {
	for _, d := range g {
		if d.Name() == name {
			return d
		}
	}
	return nil
}

func (g GoDecls) GetEnumeration(name string) *Enumeration {
	for _, d := range g {
		if e, ok := d.(*Enumeration); ok && e.name == name {
			return e
		}
	}
	return nil
}

func (g GoDecls) GetPackage() string {
	for _, d := range g {
		if _, ok := d.(*Package); ok {
			return d.Name()
		}
	}
	return ""
}

func (g GoDecls) GetStructMessage(name string) *StructMessage {
	for _, d := range g {
		if s, ok := d.(*StructMessage); ok {
			if s.name == name {
				return s
			}
		}
	}
	return nil
}

func (g GoDecls) GetTypeAlias(name string) *TypeAlias {
	for _, d := range g {
		if a, ok := d.(*TypeAlias); ok {
			if a.name == name {
				return a
			}
		}
	}
	return nil
}

func (g GoDecls) GetUnionMessage(name string) *UnionMessage {
	for _, d := range g {
		if s, ok := d.(*UnionMessage); ok {
			if s.name == name {
				return s
			}
		}
	}
	return nil
}
