package astknife

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"

	"github.com/pkg/errors"
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
func (pw *Patchwork) ParseFile(filename string, source interface{}) (*PatchworkFile, error) {
	f, err := parser.ParseFile(pw.Fset, filename, source, parser.ParseComments)
	return &PatchworkFile{Patchwork: pw, File: f}, err
}

// MustParseFile :
func (pw *Patchwork) MustParseFile(filename string, source interface{}) *PatchworkFile {
	f, err := pw.ParseFile(filename, source)
	if err != nil {
		panic(err)
	}
	return f
}

// PatchworkFile :
type PatchworkFile struct {
	*Patchwork
	File *ast.File
}

// Print :
func (pf *PatchworkFile) Print() error {
	return PrintCode(pf.Fset, pf.File)
}

// Fprint :
func (pf *PatchworkFile) Fprint(w io.Writer) error {
	return FprintCode(w, pf.Fset, pf.File)
}

// Peek :
func (pf *PatchworkFile) Peek() error {
	return PrintAST(pf.Fset, pf.File)
}

// Lookup :
func (pf *PatchworkFile) Lookup(name string) *LookupResult {
	if strings.Contains(name, ".") {
		obAndMethod := strings.SplitN(name, ".", 2)
		return LookupMethod(pf.File, obAndMethod[0], obAndMethod[1])
	}
	return LookupToplevel(pf.File, name)
}

// LookupAllMethods :
func (pf *PatchworkFile) LookupAllMethods(obname string) []*LookupResult {
	return LookupAllMethods(pf.File, obname)
}

// Append :
func (pf *PatchworkFile) Append(r *LookupResult) (ok bool, err error) {
	switch r.Type {
	case LookupTypeToplevel:
		return AppendToplevelToFile(pf.File, r.Object)
	case LookupTypeMethod:
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
