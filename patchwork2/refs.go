package patchwork2

import (
	"fmt"
	"go/ast"
)

// Ref :
type Ref struct {
	Files []*FileRef
}

// FileRef :
type FileRef struct {
	Ref      *Ref
	Decls    []declRef
	Comments []*ast.CommentGroup
	File     *ast.File
}

type (
	DeclRef interface {
		declNode()
	}

	// BadDeclRef :
	BadDeclRef struct {
		Original    *ast.BadDecl
		Replacement *ast.BadDecl
		File        *ast.File
		Object      *ast.Object
	}

	// GenDeclRef :
	GenDeclRef struct {
		Original    *ast.GenDecl
		Replacement *ast.GenDecl
		Specs       []*SpecRef
		File        *ast.File
		Object      *ast.Object
	}

	// FuncDeclRef :
	FuncDeclRef struct {
		Original    *ast.FuncDecl
		Replacement *ast.FuncDecl
		File        *ast.File
		Object      *ast.Object
	}
)

func (*BadDeclRef) declNode()  {}
func (*GenDeclRef) declNode()  {}
func (*FuncDeclRef) declNode() {}

// SpecRef :
type SpecRef struct {
	Original    ast.Spec
	Replacement ast.Spec
	File        *ast.File
	Object      *ast.Object
}

// NewRef :
func NewRef(fs []*ast.File) *Ref {
	ref := &Ref{}
	files := make([]*FileRef, len(fs))
	for i, f := range fs {
		files[i] = newFileRef(f, ref)
	}
	ref.Files = files
	return ref
}

// newFileRef :
func newFileRef(f *ast.File, ref *Ref) *FileRef {
	decls := make([]DeclRef, len(f.Decls))
	for i, decl := range f.Decls {
		decls[i] = newDeclRef(f, decl)
	}
	return &FileRef{
		Ref:      ref,
		File:     f,
		Decls:    decls,
		Comments: f.Comments,
	}
}

// newDeclRef :
func newDeclRef(f *ast.File, decl ast.Decl) DeclRef {
	switch decl := decl.(type) {
	case *ast.FuncDecl:
		return &FuncDeclRef{
			Original: decl,
			File:     f,
		}
	case *ast.GenDecl:
		specs := make([]*SpecRef, len(decl.Specs))
		for i, spec := range decl.Specs {
			specs[i] = &SpecRef{
				Original: spec,
				File:     f,
			}
		}
		return &GenDeclRef{
			Specs:    specs,
			Original: decl,
			File:     f,
		}
	case *ast.BadDecl:
		return &BadDeclRef{
			Original: decl,
			File:     f,
		}
	default:
		panic(fmt.Sprintf("invalid decl %+v\n", decl))
	}
}
