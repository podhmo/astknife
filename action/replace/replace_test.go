package replace

import (
	"bytes"
	"go/parser"
	"go/printer"
	"go/token"
	"testing"
)

func TestReplace(t *testing.T) {
	code0 := `
package p

// F : 0
func F() int {
	// this is f0's comment
	return 10
}

// G : 0
func G() int {
	// this is g0's comment
	return 20
}
`
	code1 := `
package p

// F : 1
func F() int {
	// this is f1's comment
	x := 5
	return x + x
}

// G : 1
func G() int {
	// this is g1's comment
	x := 10
	return x + x
}
`
	code2 := `
package p

// hmm

// G : 2
func G() int {
	// this is g2's comment
	x := 5
	return x + x + x + x
}
`
	// todo: doc comment is not found
	type C struct {
		msg  string
		code string
		name string
	}

	candidates := []C{
		{
			msg:  "replace f0.F to f0.F",
			code: code0,
			name: "F",
		},
		{
			msg:  "replace f0.F to f1.F",
			code: code1,
			name: "F",
		},
		{
			msg:  "replace f0.G to f1.G",
			code: code2,
			name: "G",
		},
	}

	for _, c := range candidates {
		c := c
		t.Run(c.msg, func(t *testing.T) {
			fset := token.NewFileSet()
			f0, err := parser.ParseFile(fset, "f0.go", code0, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			f1, err := parser.ParseFile(fset, "f1.go", c.code, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}

			replaced, err := Toplevel(fset, f0, f0.Scope.Lookup(c.name), f1.Scope.Lookup(c.name), f1.Comments)
			if err != nil {
				t.Fatal(err)
			}
			if !replaced {
				t.Error("must be replaced")
			}

			var b bytes.Buffer
			printer.Fprint(&b, fset, f0)
			t.Log(b.String())
		})
	}
}
