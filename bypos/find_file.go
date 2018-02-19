package bypos

import (
	"go/ast"
	"go/token"
	"sort"
)

// Sorted : sorted by pos.
type Sorted struct {
	Files []*ast.File
}

// SortFiles :
func SortFiles(files []*ast.File) Sorted {
	sort.Slice(files, func(i int, j int) bool {
		return files[i].Pos() <= files[j].Pos()
	})
	return Sorted{Files: files}
}

// FindFile :
func FindFile(sorted Sorted, pos token.Pos) *ast.File {
	var found *ast.File
	for _, f := range sorted.Files {
		if pos >= f.Pos() {
			found = f
		} else {
			return found
		}
	}
	return found
}
