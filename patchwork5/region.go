package patchwork5

import (
	"encoding/json"
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
	f        *ast.File
	Ref      Ref
	Pos      token.Pos
	End      token.Pos
	Lines    []int
	Comments []*ast.CommentGroup
}

// MarshalJSON :
func (r *Region) MarshalJSON() ([]byte, error) {
	comments := []map[string]interface{}{}
	for _, c := range r.Comments {
		comments = append(comments, map[string]interface{}{
			"text": c.Text(),
			"pos":  c.Pos(),
			"end":  c.End(),
		})
	}
	return json.Marshal(map[string]interface{}{
		"ref":      r.Ref,
		"pos":      r.Pos,
		"end":      r.End,
		"lines":    r.Lines,
		"comments": comments,
	})
}

func (r *Region) String() string {
	return fmt.Sprintf("<region %s>", r.Ref)
}

// CommentRef :
type CommentRef struct {
}

// Name :
func (r *CommentRef) Name() string {
	return ""
}

func (r *CommentRef) String() string {
	return fmt.Sprintf("<commentref>")
}

// MarshalJSON :
func (r *CommentRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"type": "commentref"})
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

// MarshalJSON :
func (r *DeclRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"name": r.name, "type": "declref"})
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

// MarshalJSON :
func (r *DeclHeadRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"type": "declhead"})
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

// MarshalJSON :
func (r *DeclTailRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"type": "decltail"})
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

// MarshalJSON :
func (r *SpecRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"type": "specref", "name": r.name})
}
