package patchwork4

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/podhmo/astknife/patchwork4/lookup"
)

var errNotfound = errors.New("not found")

// New :
func New(fset *token.FileSet, f *ast.File) *Patchwork {
	return &Patchwork{
		Fset:      fset,
		File:      f,
		Replacing: map[ast.Node]*lookup.Result{},
	}
}

// Patchwork :
type Patchwork struct {
	Fset      *token.FileSet
	File      *ast.File
	Replacing map[ast.Node]*lookup.Result
	Appending []*lookup.Result
	// todo: support Removings
}

// Append :
func (p *Patchwork) Append(r *lookup.Result) error {
	if r == nil {
		return errNotfound
	}
	p.Appending = append(p.Appending, r)
	return nil
}

// Replace :
func (p *Patchwork) Replace(pat *lookup.Result, rep *lookup.Result) error {
	if rep == nil {
		return errNotfound
	}
	if pat == nil {
		return errNotfound
	}

	if pat.Object.Kind != rep.Object.Kind {
		return &Unsupported{pat, fmt.Sprintf("conflict kind %s != %s", pat.Object.Kind, rep.Object.Kind)}
	}
	node, ok := pat.Object.Decl.(ast.Node)
	if !ok {
		return &Unsupported{pat, "not ast.Node"}
	}
	p.Replacing[node] = rep
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
