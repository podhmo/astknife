package patchwork

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// Patchwork : (todo rename)
type Patchwork struct {
	Fset  *token.FileSet
	scope *scope
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
	if pw.scope == nil {
		pw.scope = newscope()
	}
	return pw
}

// WithFileSet :
func WithFileSet(fset *token.FileSet) func(*Patchwork) {
	return func(pw *Patchwork) {
		pw.Fset = fset
	}
}

// ParseFile :
func (pw *Patchwork) ParseFile(filename string, source interface{}) (*File, error) {
	file, err := parser.ParseFile(pw.Fset, filename, source, parser.ParseComments)
	pw.scope.AddFile(filename, file)
	f := &File{Patchwork: pw, File: file}
	return f, err
}

// ParseAST :
func (pw *Patchwork) ParseAST(filename string, file *ast.File) (*File, error) {
	pw.scope.AddFile(filename, file)
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
