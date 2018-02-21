package lookup

import (
	"go/ast"
)

// Type :
type Type string

const (
	// TypeToplevel : toplevel (e.g. toplevel function, struct definition)
	TypeToplevel = Type("toplevel")
	// TypeMethod : method (e.g. method function, struct definition)
	TypeMethod = Type("method")
)

// Result :
type Result struct {
	Type     Type
	FuncDecl *ast.FuncDecl
	Object   *ast.Object
}

// Name :
func (r *Result) Name() string {
	switch r.Type {
	case TypeToplevel:
		return r.Object.Name
	case TypeMethod:
		return r.FuncDecl.Name.Name
	}
	return "<nil>"
}
