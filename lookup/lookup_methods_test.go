package lookup

import (
	"fmt"
	"go/parser"
	"go/token"
	"testing"
)

func TestMethods(t *testing.T) {
	source := `
package p

func Hello() string {
	return "hello"
}

type S struct {}
func (s *S) Hello() string {
	return "s.hello"
}
func (s *S) String() string {
	return "s"
}

type S2 struct {}
func (s *S2) Hello() string {
	return "s2.hello"
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("method", func(t *testing.T) {
		candidats := []struct {
			obname   string
			name     string
			notfound bool
		}{
			{
				obname: "S",
				name:   "Hello",
			},
			{
				obname: "S",
				name:   "String",
			},
			{
				obname:   "S",
				name:     "NotFound",
				notfound: true,
			},
		}
		for _, c := range candidats {
			c := c
			t.Run(fmt.Sprintf("lookup %s.%s'", c.obname, c.name), func(t *testing.T) {
				got := Method(f, c.obname, c.name)
				if c.notfound {
					if got != nil {
						t.Fatalf("should %s is not found, but found %s", c.name, got.Name())
					}
					return
				}

				if got == nil {
					t.Fatalf("should %s is found, but not found", c.name)
				}

				if got.Name() != c.name {
					t.Fatalf("should method name is %s, but got %s", c.name, got.Name())
				}
			})
		}
	})

	t.Run("allmethods", func(t *testing.T) {
		candidats := []struct {
			obname        string
			expectedCount int
		}{
			{
				obname:        "S",
				expectedCount: 2,
			},
			{
				obname:        "S2",
				expectedCount: 1,
			},
			{
				obname:        "S3",
				expectedCount: 0,
			},
		}

		for _, c := range candidats {
			c := c
			t.Run(fmt.Sprintf("%s's methods", c.obname), func(t *testing.T) {
				methods := AllMethods(f, c.obname)
				if len(methods) != c.expectedCount {
					t.Fatalf("should len(methods) == %d, but got %d", c.expectedCount, len(methods))
				}
			})
		}
	})
}
