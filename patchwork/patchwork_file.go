package patchwork

import (
	"go/ast"
	"io"

	"github.com/podhmo/astknife/action"
	"github.com/podhmo/astknife/lookup"
	"github.com/podhmo/astknife/printer"
)

// File :
type File struct {
	*Patchwork
	File *ast.File
}

// FprintCode :
func (pf *File) FprintCode(w io.Writer) error {
	return printer.FprintCode(w, pf.Fset, pf.File)
}

// PrintCode :
func (pf *File) PrintCode() error {
	return printer.PrintCode(pf.Fset, pf.File)
}

// PrintAST :
func (pf *File) PrintAST() error {
	return printer.PrintAST(pf.Fset, pf.File)
}

// Lookup :
func (pf *File) Lookup(name string) *lookup.Result {
	return pf.lookup.Lookup(name)
}

// LookupAllMethods :
func (pf *File) LookupAllMethods(obname string) []*lookup.Result {
	// todo: xxx
	return pf.lookup.AllMethods(obname)
}

// Append :
func (pf *File) Append(r *lookup.Result) (ok bool, err error) {
	return action.Append(pf.lookup, pf.File, r)
}

// Replace :
func (pf *File) Replace(r *lookup.Result) (ok bool, err error) {
	return action.Replace(pf.lookup, pf.File, r)
}

// AppendOrReplace : upsert
func (pf *File) AppendOrReplace(r *lookup.Result) (ok bool, err error) {
	return action.AppendOrReplace(pf.lookup, pf.File, r)
}

// Wrap : xxx
func (pf *File) Wrap(pw *Patchwork) *File {
	return &File{
		Patchwork: &Patchwork{
			Fset:   pw.Fset,
			lookup: pw.lookup.With(pf.File),
		},
		File: pf.File,
	}
}
