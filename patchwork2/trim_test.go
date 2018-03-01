package patchwork2

import (
	"fmt"
	"go/parser"
	"go/token"
	"testing"
)

func TestCalcSize(t *testing.T) {
	source := `
package p
// F :
func F() int {
	return 1 + 1
} // xxx
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	size := fset.File(f.Pos()).Size()
	base := fset.File(f.Pos()).Base()
	fmt.Println(base, size)
	{
		fmt.Println("----------------------------------------")
		size := trim(fset, NewRef().NewFileRef(f))
		fmt.Println("----------------------------------------")
		fmt.Println(size)
	}
}
