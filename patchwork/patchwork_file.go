package patchwork

import (
	"go/ast"
	"io"
	"strings"

	"github.com/podhmo/astknife/action"
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
			return lookup.MethodByObject(f, ob.Object, name)
		})
	}
	return pf.scope.Lookup(name, pf.File, lookup.Toplevel)
}

// LookupAllMethods :
func (pf *File) LookupAllMethods(obname string) []*lookup.Result {
	// todo: xxx
	return lookup.AllMethods(pf.File, obname)
}

// Append :
func (pf *File) Append(r *lookup.Result) (ok bool, err error) {
	return action.Append(pf.File, r)
}

// Replace :
func (pf *File) Replace(r *lookup.Result) (ok bool, err error) {
	return action.Replace(pf.File, r)
}

// AppendOrReplace : upsert
func (pf *File) AppendOrReplace(r *lookup.Result) (ok bool, err error) {
	return action.AppendOrReplace(pf.File, r)
}
