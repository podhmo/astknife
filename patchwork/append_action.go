package patchwork

import (
	"go/ast"
	"go/token"

	"github.com/pkg/errors"
	"github.com/podhmo/astknife/lookup"
)

// Append :
func (pf *File) Append(r *lookup.Result) (ok bool, err error) {
	if r == nil {
		return false, ErrNotFound
	}

	switch r.Type {
	case lookup.TypeToplevel:
		return AppendToplevelToFile(pf.File, r.Object)
	case lookup.TypeMethod:
		return AppendFunctionToFile(pf.File, r.FuncDecl)
	default:
		return false, errors.New("not implemented")
	}
}

// todo: all object types support
// // The list of possible Object kinds.
// const (
// 	Bad ObjKind = iota // for error handling
// 	Pkg                // package
// 	Con                // constant
// 	Typ                // type
// 	Var                // variable
// 	Fun                // function or method
// 	Lbl                // label
// )

// todo: comment support

// AppendToplevelToFile :
func AppendToplevelToFile(dst *ast.File, ob *ast.Object) (ok bool, err error) {
	if ob == nil {
		return
	}

	if ob := dst.Scope.Lookup(ob.Name); ob != nil {
		err = errors.Errorf("%s is already existed, in scope", ob.Name)
		return
	}
	dst.Scope.Insert(ob)

	switch ob.Kind {
	case ast.Typ:
		dst.Decls = append(dst.Decls, &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{ob.Decl.(ast.Spec)},
		})
		ok = true
	case ast.Fun:
		if decl, can := ob.Decl.(*ast.FuncDecl); can {
			return AppendFunctionToFile(dst, decl)
		}
		err = errors.Errorf("unsupported object type %s (kind=%q)", ob.Type, ob.Kind)
		return
	}
	return
}

// AppendFunctionToFile :
func AppendFunctionToFile(dst *ast.File, decl *ast.FuncDecl) (ok bool, err error) {
	if decl == nil {
		return
	}

	dst.Decls = append(dst.Decls, decl)
	ok = true
	return
}
