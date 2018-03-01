package patchwork2

import (
	"go/ast"

	"github.com/podhmo/astknife/patchwork2/lookup"
)

// Ref :
type Ref struct {
	Files []*FileRef
}

// FileRef :
type FileRef struct {
	Ref      *Ref
	Decls    []*DeclRef
	Comments []*ast.CommentGroup
	File     *ast.File
}

// DeclRef :
type DeclRef struct {
	Original    ast.Decl
	Replacement ast.Decl
	Specs       []*SpecRef
	File        *ast.File
	Result      *lookup.Result
	Comments    []*ast.CommentGroup
}

// SpecRef :
type SpecRef struct {
	Original    ast.Spec
	Replacement ast.Spec
	File        *ast.File
	Result      *lookup.Result
	Comments    []*ast.CommentGroup
}

// NewRef :
func NewRef(fs ...*ast.File) *Ref {
	ref := &Ref{}
	files := make([]*FileRef, len(fs))
	for i, f := range fs {
		files[i] = ref.NewFileRef(f)
	}
	ref.Files = files
	return ref
}

// NewFileRef :
func (r *Ref) NewFileRef(f *ast.File) *FileRef {
	fileref := &FileRef{
		Ref:      r,
		File:     f,
		Comments: f.Comments,
	}
	decls := make([]*DeclRef, len(f.Decls))
	for i, decl := range f.Decls {
		decls[i] = newDeclRef(f, decl)
	}
	fileref.Decls = decls
	return fileref
}

// newDeclRef :
func newDeclRef(f *ast.File, decl ast.Decl) *DeclRef {
	declref := &DeclRef{
		Original: decl,
		File:     f,
	}
	if t, ok := decl.(*ast.GenDecl); ok {
		specs := make([]*SpecRef, len(t.Specs))
		for i, spec := range t.Specs {
			specs[i] = &SpecRef{
				Original: spec,
				File:     f,
			}
		}
		declref.Specs = specs
	}
	return declref
}
