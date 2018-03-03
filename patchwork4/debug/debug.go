package debug

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// Debug :
type Debug struct {
	Fset      *token.FileSet
	SourceMap map[string]string
}

// ParseSource :
func (d *Debug) ParseSource(filename string, source string) (*ast.File, error) {
	f, err := parser.ParseFile(d.Fset, filename, source, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	k := d.Fset.File(f.Pos()).Name()
	d.SourceMap[k] = source
	return f, nil
}

// New :
func New(fset *token.FileSet) *Debug {
	return &Debug{
		Fset:      fset,
		SourceMap: map[string]string{},
	}
}
