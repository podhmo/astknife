package patchwork3

import (
	"errors"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"io"

	"github.com/podhmo/astknife/patchwork3/lookup"
)

var errNotfound = errors.New("not found")

// Patchwork :
type Patchwork struct {
	Fset    *token.FileSet
	Printer *printer.Config
}

// New :
func New(fset *token.FileSet) *Patchwork {
	return &Patchwork{Fset: fset, Printer: &printer.Config{Tabwidth: 8}}
}

// NewArea :
func (w *Patchwork) NewArea(f *ast.File) Area {
	return &SingleArea{
		File: f,
		Node: f,
		p:    w,
	}
}

// Append :
func (w *Patchwork) Append(area Area, r *lookup.Result) (Area, error) {
	return area.Append(r)
}

// Area :
type Area interface {
	Display(w io.Writer) error
	Append(r *lookup.Result) (Area, error)
}

// SingleArea :
type SingleArea struct {
	p        *Patchwork
	Node     ast.Node
	File     *ast.File
	Comments []*ast.CommentGroup
}

// Display :
func (s *SingleArea) Display(w io.Writer) error {
	if len(s.Comments) > 0 {
		target := &printer.CommentedNode{Node: s.Node, Comments: s.Comments}
		return s.p.Printer.Fprint(w, s.p.Fset, target)
	}
	return s.p.Printer.Fprint(w, s.p.Fset, s.Node)
}

// Append :
func (s *SingleArea) Append(r *lookup.Result) (Area, error) {
	if r == nil {
		return nil, errNotfound
	}
	node := r.Object.Decl.(ast.Node)
	return &MixedArea{
		p: s.p,
		Areas: []Area{s, &SingleArea{
			File:     r.File,
			Node:     node,
			Comments: ast.NewCommentMap(s.p.Fset, node, r.File.Comments).Filter(node).Comments(),
			p:        s.p,
		}},
	}, nil
}

// MixedArea :
type MixedArea struct {
	p     *Patchwork
	Areas []Area
}

// Append :
func (m *MixedArea) Append(r *lookup.Result) (Area, error) {
	if r == nil {
		return nil, errNotfound
	}
	return &MixedArea{
		p: m.p,
		Areas: append(m.Areas, &SingleArea{
			File: r.File,
			Node: r.Object.Decl.(ast.Node),
			p:    m.p,
		}),
	}, nil
}

// Display :
func (m *MixedArea) Display(w io.Writer) error {
	for _, area := range m.Areas {
		if err := area.Display(w); err != nil {
			return err
		}
	}
	return nil
}

// Unsupported :
type Unsupported struct {
	Result *lookup.Result
	Reason string
}

// Unsupported :
func (e *Unsupported) Error() string {
	return fmt.Sprintf("unsupported object %s (reason=%s, type=%s, kind=%q)",
		e.Result.String(), e.Reason, e.Result.Object.Type, e.Result.Object.Kind,
	)
}
