package patchwork

import (
	"bytes"
	"fmt"
	"testing"
)

// TestAppend
func TestAppend(t *testing.T) {
	source := `
package p
type S struct {}
`
	type C struct {
		source  string
		source2 string
		name    string
		hasErr  bool
	}

	candidates := []C{
		{
			source: source,
			name:   "S2",
			source2: `
		package p
		type S2 struct {}
		`,
		},
		{
			source:  source,
			name:    "S",
			source2: source,
			hasErr:  true,
		},
		{
			source: source,
			name:   "Hello",
			source2: `
package p
func Hello() string {
	return "hello"
}
type S2 struct {}
func (s *S2) Hello() string {
	return "s.hello"
}
`,
		},
		{
			source: source,
			name:   "S.Hello",
			source2: `
package p
func (s *S) Hello() string {
	return "s.hello"
}
`,
		},
	}

	for _, c := range candidates {
		c := c
		t.Run(fmt.Sprintf("append %s", c.name), func(t *testing.T) {
			pf := NewPatchwork().MustParseFile("f0", source)
			pf2 := pf.MustParseFile("f1", c.source2)

			t.Logf("input (%s)\n%s\n", c.name, c.source)
			t.Logf("append (%s)\n%s\n", c.name, c.source2)

			ok, err := pf.Append(pf2.Lookup(c.name))

			if c.hasErr {
				t.Logf("should error %s", err)
				if err == nil {
					t.Fatal("error is expected, but no error")
				}
				return
			}

			var b bytes.Buffer
			if err := pf.Fprint(&b); err != nil {
				t.Fatal(err)
			}
			t.Logf("output\n%s\n", b.String())

			if err != nil {
				t.Fatal(err)
			}

			if !ok {
				t.Fatal("must appended")
			}

			if pf.Lookup(c.name) == nil {
				t.Fatalf("cannot lookup appended object (%q)", c.name)
			}
		})
	}
}
