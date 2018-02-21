package action

import (
	"go/ast"

	"github.com/pkg/errors"
	"github.com/podhmo/astknife/action/append"
	"github.com/podhmo/astknife/action/replace"
	"github.com/podhmo/astknife/lookup"
)

// Append :
func Append(f *ast.File, r *lookup.Result) (ok bool, err error) {
	if r == nil {
		return false, ErrNoEffect
	}

	switch r.Type {
	case lookup.TypeToplevel:
		return append.ToplevelToFile(f, r.Object)
	case lookup.TypeMethod:
		return append.FunctionToFile(f, r.FuncDecl)
	default:
		return false, errors.New("not implemented")
	}
}

// Replace :
func Replace(f *ast.File, r *lookup.Result) (ok bool, err error) {
	if r == nil {
		return false, ErrNoEffect
	}

	switch r.Type {
	case lookup.TypeToplevel:
		drObject := f.Scope.Lookup(r.Name())
		if drObject == nil {
			return false, ErrNotFound
		}
		return replace.ToplevelToFile(f, drObject, r.Object)
	case lookup.TypeMethod:
		dr := lookup.MethodByObject(f, r.Object, r.Name())
		if dr == nil {
			return false, ErrNotFound
		}
		return replace.MethodToFile(f, r.Object, dr.FuncDecl, r.FuncDecl)
	default:
		return false, errors.New("not implemented")
	}
}

// AppendOrReplace : upsert
func AppendOrReplace(f *ast.File, r *lookup.Result) (ok bool, err error) {
	if r == nil {
		return false, ErrNoEffect
	}

	switch r.Type {
	case lookup.TypeToplevel:
		drObject := f.Scope.Lookup(r.Name())
		if drObject == nil {
			return append.ToplevelToFile(f, r.Object)
		}
		return replace.ToplevelToFile(f, drObject, r.Object)
	case lookup.TypeMethod:
		dr := lookup.MethodByObject(f, r.Object, r.Name())
		if dr == nil {
			return append.FunctionToFile(f, r.FuncDecl)
		}
		return replace.MethodToFile(f, r.Object, dr.FuncDecl, r.FuncDecl)
	default:
		return false, errors.New("not implemented")
	}
}
