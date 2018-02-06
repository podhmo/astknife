package astknife

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"
)

func TestFindComments(t *testing.T) {
	fset := token.NewFileSet()

	source := `
package p

// S : *this is S*
type S struct {
	// Name : *name of S*
	Name string
	Age int // *this is age*
}

// String : *for stringer*
func (s *S) String() string {
	return s.Name
}
`
	conf := &types.Config{
		Error: func(err error) {
			t.Fatalf("error when typecheck %s", err)
		},
	}
	file, _ := parser.ParseFile(fset, "", source, parser.ParseComments)
	files := []*ast.File{file}
	pkg, err := conf.Check("p", fset, files, nil)
	if err != nil {
		t.Fatalf("error when typecheck2 %s", err)
	}

	type C struct {
		msg     string
		getPos  func() token.Pos
		comment string
	}
	candidates := []C{
		{
			msg: "toplevel struct comments",
			getPos: func() token.Pos {
				// todo: fix asjutment -1
				return pkg.Scope().Lookup("S").Pos() - 1 // xxx
			},
			comment: "*this is S*",
		},
		{
			msg: "toplevel struct, field comments",
			getPos: func() token.Pos {
				ob := pkg.Scope().Lookup("S")
				internal := ob.Type().Underlying().(*types.Struct)
				return internal.Field(0).Pos()
			},
			comment: "*name of S*",
		},
		{
			msg: "toplevel struct, field comments, end of line",
			getPos: func() token.Pos {
				ob := pkg.Scope().Lookup("S")
				internal := ob.Type().Underlying().(*types.Struct)
				return internal.Field(1).Pos()
			},
			comment: "*this is age*",
		},
		{
			msg: "toplevel struct, method definition comments",
			getPos: func() token.Pos {
				ob := pkg.Scope().Lookup("S")
				method, _, _ := types.LookupFieldOrMethod(ob.Type(), true, pkg, "String")
				return method.Pos() - 1 // xxx:
			},
			comment: "*for stringer*",
		},
		// todo: const
	}

	for _, c := range candidates {
		c := c
		t.Run(c.msg, func(t *testing.T) {
			comments := FindCommentsByPos(files, c.getPos())
			if !strings.Contains(comments.Text(), c.comment) {
				t.Errorf("expected contains %q, but got %q", c.comment, comments.Text())
			}
		})
	}
}
