package lookup

import (
	"fmt"
	"go/ast"
)

// Type :
type Type string

const (
	// TypeToplevel :
	TypeToplevel = Type("toplevel")
	// TypeMethod :
	TypeMethod = Type("method")
)

// Result :
type Result struct {
	File     *ast.File
	Type     Type
	Recv     string
	Name     string
	FuncDecl *ast.FuncDecl
	Object   *ast.Object
}

// String
func (r *Result) String() string {
	switch r.Type {
	case TypeToplevel:
		return r.Name
	case TypeMethod:
		return fmt.Sprintf("%s.%s", r.Recv, r.Name)
	}
	return "<nil>"
}
