package patchwork

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/astknife/lookup"
	"github.com/podhmo/astknife/printer"
)

// Patchwork : (todo rename)
type Patchwork struct {
	Fset *token.FileSet
}

// NewPatchwork :
func NewPatchwork() *Patchwork {
	return &Patchwork{Fset: token.NewFileSet()}
}

// ParseFile :
func (pw *Patchwork) ParseFile(filename string, source interface{}) (*File, error) {
	f, err := parser.ParseFile(pw.Fset, filename, source, parser.ParseComments)
	return &File{Patchwork: pw, File: f}, err
}

// MustParseFile :
func (pw *Patchwork) MustParseFile(filename string, source interface{}) *File {
	f, err := pw.ParseFile(filename, source)
	if err != nil {
		panic(err)
	}
	return f
}

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
		return lookup.LookupMethod(pf.File, obAndMethod[0], obAndMethod[1])
	}
	return lookup.LookupToplevel(pf.File, name)
}

// LookupAllMethods :
func (pf *File) LookupAllMethods(obname string) []*lookup.Result {
	return lookup.LookupAllMethods(pf.File, obname)
}

// Append :
func (pf *File) Append(r *lookup.Result) (ok bool, err error) {
	switch r.Type {
	case lookup.TypeToplevel:
		return AppendToplevelToFile(pf.File, r.Object)
	case lookup.TypeMethod:
		// toDO
		return false, errors.New("not implemented")
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
			dst.Decls = append(dst.Decls, decl)
			ok = true
			return
		}
		err = errors.Errorf("unsupported object type %s (kind=%q)", ob.Type, ob.Kind)
		return
	}
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
