package patchwork

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/podhmo/astknife/lookup"
)

// Patchwork : (todo rename)
type Patchwork struct {
	Fset   *token.FileSet
	lookup *lookup.Lookup
}

// NewPatchwork :
func NewPatchwork(opts ...func(*Patchwork)) *Patchwork {
	pw := &Patchwork{}
	for _, op := range opts {
		op(pw)
	}
	if pw.Fset == nil {
		pw.Fset = token.NewFileSet()
	}
	if pw.lookup == nil {
		pw.lookup = lookup.New()
	}
	return pw
}

// WithFileSet :
func WithFileSet(fset *token.FileSet) func(*Patchwork) {
	return func(pw *Patchwork) {
		pw.Fset = fset
	}
}

// WithLookup :
func WithLookup(lookup *lookup.Lookup) func(*Patchwork) {
	return func(pw *Patchwork) {
		pw.lookup = lookup
	}
}

// ParseFile :
func (pw *Patchwork) ParseFile(filename string, source interface{}) (*File, error) {
	file, err := parser.ParseFile(pw.Fset, filename, source, parser.ParseComments)
	pw.lookup.Files = append(pw.lookup.Files, file)
	f := &File{Patchwork: pw, File: file}
	return f, err
}

// ParseAST :
func (pw *Patchwork) ParseAST(filename string, file *ast.File) (*File, error) {
	pw.lookup.Files = append(pw.lookup.Files, file)
	f := &File{Patchwork: pw, File: file}
	return f, nil
}

// MustParseFile :
func (pw *Patchwork) MustParseFile(filename string, source interface{}) *File {
	f, err := pw.ParseFile(filename, source)
	if err != nil {
		panic(err)
	}
	return f
}
