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

// F :
func F() int {
	return 10
}

// G :
func G() int {
	return 20
}
`
	code1 := `
package p

// F :
func F() int {
	x := 5
	return x + x
}
`
	code2 := `
package p

// G :
func G() int {
	x := 10
	return x + x
}
`

	type C struct {
		msg  string
		code string
		name string
	}

	candidates := []C{
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

			replaced, err := ToplevelToFile(f0, f0.Scope.Lookup(c.name), f1.Scope.Lookup(c.name))
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
