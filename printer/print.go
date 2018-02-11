package printer

import (
	"go/ast"
	"go/printer"
	"go/token"
	"io"
	"os"
)

// FprintCode :
func FprintCode(w io.Writer, fset *token.FileSet, node ast.Node) error {
	config := &printer.Config{Tabwidth: 8, Mode: printer.TabIndent}
	return config.Fprint(w, fset, node)
}

// PrintCode :
func PrintCode(fset *token.FileSet, node ast.Node) error {
	config := &printer.Config{Tabwidth: 8, Mode: printer.TabIndent}
	return config.Fprint(os.Stdout, fset, node)
}

// PrintAST :
func PrintAST(fset *token.FileSet, node ast.Node) error {
	return ast.Print(fset, node)
}
