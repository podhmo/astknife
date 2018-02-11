package patchwork

import (
	"go/ast"
	"io"
	"strings"

	"github.com/podhmo/astknife/lookup"
	"github.com/podhmo/astknife/printer"
)

// File :
type File struct {
	*Patchwork
	File *ast.File
}

// FprintCode :
func (pf *File) FprintCode(w io.Writer) error {
	return printer.FprintCode(w, pf.Fset, pf.File)
}

// PrintCode :
func (pf *File) PrintCode() error {
	return printer.PrintCode(pf.Fset, pf.File)
}

// PrintAST :
func (pf *File) PrintAST() error {
	return printer.PrintAST(pf.Fset, pf.File)
}

// Lookup :
func (pf *File) Lookup(name string) *lookup.Result {
	if strings.Contains(name, ".") {
		obAndMethod := strings.SplitN(name, ".", 2)
		ob := pf.scope.Lookup(obAndMethod[0], pf.File, lookup.Toplevel)
		if ob == nil {
			return nil
		}
		return pf.scope.Lookup(obAndMethod[1], pf.File, func(f *ast.File, name string) *lookup.Result {
			return lookup.MethodFromObject(f, ob.Object, name)
		})
	}
	return pf.scope.Lookup(name, pf.File, lookup.Toplevel)
}

// LookupAllMethods :
func (pf *File) LookupAllMethods(obname string) []*lookup.Result {
	return lookup.AllMethods(pf.File, obname)
}