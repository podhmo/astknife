package patchwork4

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"testing"

	"github.com/podhmo/astknife/patchwork4/debug"
	"github.com/podhmo/astknife/patchwork4/lookup"
	"github.com/podhmo/printer"
)

func TestReplace(t *testing.T) {
	code0 := `
package p

// F : 0
func F() int {
    // this is f0's comment
    return 10
}

// toplevel comment

// S : 0
type S struct {
    // Name :
    Name string // name
    // Age :
    Age string // age
    // Nickname :
    Nickname string // nickname
}

// G : 0
func G() int {
	// this is g0's comment
	return 10 + 10
}

// H : 0
func H() int {
	// this is h0's comment
	return 10 + 10
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

// S : 2
type S struct {
    // Name :
    Name string // name
    // Age :
    Age string // age
    // Nickname :
    Nickname string // nickname
}

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
		// {
		// 	msg:  "replace f0.F to f0.F",
		// 	code: code0,
		// 	name: "F",
		// },
		// {
		// 	msg:  "replace f0.F to f1.F",
		// 	code: code1,
		// 	name: "F",
		// },
		// {
		// 	msg:  "replace f0.G to f1.G",
		// 	code: code2,
		// 	name: "G",
		// },
		{
			msg:  "replace f0.S to f1.S",
			code: code2,
			name: "S",
		},
	}
	_ = code1
	_ = code2
	for _, c := range candidates {
		c := c
		t.Run(c.msg, func(t *testing.T) {
			fset := token.NewFileSet()
			debug := debug.New(fset)
			f0, err := debug.ParseSource("f0.go", code0)
			if err != nil {
				t.Fatal(err)
			}
			f1, err := debug.ParseSource("f1.go", c.code)
			if err != nil {
				t.Fatal(err)
			}
			p := New(fset, f0, WithDebug(debug))
			if err := p.Replace(lookup.Lookup(c.name, f0), lookup.Lookup(c.name, f1)); err != nil {
				t.Fatal(err)
			}

			var b bytes.Buffer
			got := ToFile(p, "newf0.go")
			printer.Fprint(&b, fset, got)

			// tf := fset.File(got.Pos())
			// ast.Inspect(got, func(node ast.Node) bool {
			// 	if node != nil {
			// 		fmt.Printf("%02d: %T (%d,%d)\n", tf.Line(node.Pos()), node, int(node.Pos())-tf.Base(), int(node.End())-tf.Base())
			// 	}
			// 	return true
			// })
			t.Log(b.String())
		})
	}
}

func dumpPositions(f *ast.File) {
	ast.Inspect(f, func(node ast.Node) bool {
		if node != nil {
			if _, ok := node.(ast.Decl); ok {
				fmt.Println("-")
			}
			switch x := node.(type) {
			case *ast.Ident:
				fmt.Printf("%T(%s) %v-%v *size=%v* @ %v-%v\n", x, x.Name, x.Pos(), x.End(), x.End()-x.Pos(), x.Pos()-f.Pos(), x.End()-f.Pos())
			case *ast.CommentGroup:
				fmt.Printf("%T(%q) %v-%v *size=%v* @ %v-%v\n", x, x.Text(), x.Pos(), x.End(), x.End()-x.Pos(), x.Pos()-f.Pos(), x.End()-f.Pos())
			default:
				fmt.Printf("%T %v-%v *size=%v* @ %v-%v\n", node, node.Pos(), node.End(), node.End()-node.Pos(), node.Pos()-f.Pos(), node.End()-f.Pos())
			}
		}
		return true
	})
}
