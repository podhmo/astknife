package patchwork

import (
	"go/parser"
	"go/token"
)

// Patchwork : (todo rename)
type Patchwork struct {
	Fset  *token.FileSet
	scope *scope
}

// NewPatchwork :
func NewPatchwork() *Patchwork {
	return &Patchwork{Fset: token.NewFileSet(), scope: newscope()}
}

// ParseFile :
func (pw *Patchwork) ParseFile(filename string, source interface{}) (*File, error) {
	file, err := parser.ParseFile(pw.Fset, filename, source, parser.ParseComments)
	pw.scope.AddFile(filename, file)
	f := &File{Patchwork: pw, File: file}
	return f, err
}

// MustParseFile :
func (pw *Patchwork) MustParseFile(filename string, source interface{}) *File {
	f, err := pw.ParseFile(filename, source)
	if err != nil {
		panic(err)
	}
	return f
}
