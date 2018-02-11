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
	return printer.Fprint(w, fset, node)
}

// PrintCode :
func PrintCode(fset *token.FileSet, node ast.Node) error {
	return printer.Fprint(os.Stdout, fset, node)
}

// PrintAST :
func PrintAST(fset *token.FileSet, node ast.Node) error {
	return ast.Print(fset, node)
}
