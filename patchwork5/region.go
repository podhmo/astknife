package patchwork5

import (
	"fmt"
	"go/ast"
	"go/token"
)

// File :
type File struct {
	Regions []*Region
}

// NewFile :
func NewFile(fset *token.FileSet, f *ast.File) *File {
	if f == nil {
		return nil
	}
	return parseASTFile(fset, f)
}

// Ref :
type Ref interface {
	Name() string
	fmt.Stringer
}

// Region :
type Region struct {
	f      *ast.File
	Ref    Ref
	Origin int
	Lines  []int
}

func (r *Region) String() string {
	return fmt.Sprintf("<region %s>", r.Ref)
}

// CommentRef :
type CommentRef struct {
	Comments []*ast.CommentGroup
}

// Name :
func (r *CommentRef) Name() string {
	return ""
}

func (r *CommentRef) String() string {
	return fmt.Sprintf("<commentref>")
}

// DeclRef :
type DeclRef struct {
	name string
	Decl ast.Decl
}

// Name :
func (r *DeclRef) Name() string {
	return r.name
}

func (r *DeclRef) String() string {
	return fmt.Sprintf("<declref name=%q>", r.name)
}

// DeclHeadRef :
type DeclHeadRef struct {
	decl *ast.GenDecl
}

// Name :
func (r *DeclHeadRef) Name() string {
	return ""
}

func (r *DeclHeadRef) String() string {
	return fmt.Sprintf("<declhead>")
}

// DeclTailRef :
type DeclTailRef struct {
	decl *ast.GenDecl
}

// Name :
func (r *DeclTailRef) Name() string {
	return ""
}

func (r *DeclTailRef) String() string {
	return fmt.Sprintf("<decltail>")
}

// SpecRef :
type SpecRef struct {
	name string
	Spec ast.Spec
}

// Name :
func (r *SpecRef) Name() string {
	return r.name
}

func (r *SpecRef) String() string {
	return fmt.Sprintf("<specref name=%s>", r.name)
}
