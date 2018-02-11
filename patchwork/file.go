package patchwork

import (
	"go/ast"
	"io"
	"strings"

	"github.com/pkg/errors"
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
	return lookup.AllMethods(pf.File, obname)
}

// Append :
func (pf *File) Append(r *lookup.Result) (ok bool, err error) {
	if r == nil {
		return false, ErrNoEffect
	}

	switch r.Type {
	case lookup.TypeToplevel:
		return appendToplevelToFile(pf.File, r.Object)
	case lookup.TypeMethod:
		return appendFunctionToFile(pf.File, r.FuncDecl)
	default:
		return false, errors.New("not implemented")
	}
}

// Replace :
func (pf *File) Replace(r *lookup.Result) (ok bool, err error) {
	if r == nil {
		return false, ErrNoEffect
	}

	switch r.Type {
	case lookup.TypeToplevel:
		drObject := pf.File.Scope.Lookup(r.Name())
		if drObject == nil {
			err = errors.Errorf("%s is not existed, in scope", r.Name())
			return
		}
		return replaceToplevelToFile(pf.File, drObject, r.Object)
	case lookup.TypeMethod:
		dr := pf.scope.Lookup(r.Name(), pf.File, func(f *ast.File, name string) *lookup.Result {
			return lookup.MethodByObject(f, r.Object, name)
		})
		if dr == nil {
			return false, ErrNoEffect
		}
		return replaceMethodToFile(pf.File, r.Object, dr.FuncDecl, r.FuncDecl)
	default:
		return false, errors.New("not implemented")
	}
}

// AppendOrReplace : upsert
func (pf *File) AppendOrReplace(r *lookup.Result) (ok bool, err error) {
	if r == nil {
		return false, ErrNoEffect
	}

	switch r.Type {
	case lookup.TypeToplevel:
		drObject := pf.File.Scope.Lookup(r.Name())
		if drObject == nil {
			return appendToplevelToFile(pf.File, r.Object)
		}
		return replaceToplevelToFile(pf.File, drObject, r.Object)
	case lookup.TypeMethod:
		dr := pf.scope.Lookup(r.Name(), pf.File, func(f *ast.File, name string) *lookup.Result {
			return lookup.MethodByObject(f, r.Object, name)
		})
		if dr == nil {
			return appendFunctionToFile(pf.File, r.FuncDecl)
		}
		return replaceMethodToFile(pf.File, r.Object, dr.FuncDecl, r.FuncDecl)
	default:
		return false, errors.New("not implemented")
	}
}
