package astknife

import (
	"go/ast"
	"go/token"
)

// FindFileByPos :
func FindFileByPos(files []*ast.File, pos token.Pos) *ast.File {
	var found *ast.File
	for _, f := range files {
		if pos >= f.Pos() {
			found = f
		} else {
			return found
		}
	}
	return found
}
