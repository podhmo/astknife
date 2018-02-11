package patchwork

import (
	"go/ast"
	"go/token"
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

// Print :
func (pf *File) Print() error {
	return printer.PrintCode(pf.Fset, pf.File)
}

// Fprint :
func (pf *File) Fprint(w io.Writer) error {
	return printer.FprintCode(w, pf.Fset, pf.File)
}

// Peek :
func (pf *File) Peek() error {
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

var (
	// ErrNotFound :
	ErrNotFound = errors.New("not found")
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
