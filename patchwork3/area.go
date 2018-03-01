package patchwork3

import (
	"go/ast"
	"go/printer"
	"go/token"
	"io"
)

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

// Area :
type Area interface {
	Display(w io.Writer) error
}

// SingleArea :
type SingleArea struct {
	p    *Patchwork
	Node ast.Node
	File *ast.File
}

// Display :
func (s *SingleArea) Display(w io.Writer) error {
	return s.p.Printer.Fprint(w, s.p.Fset, s.Node)
}

// // MixedArea :
// type MixedArea struct {
// 	Areas []Area
// }
