package patchwork

import (
	"bytes"
	"strings"
	"testing"
)

// TestReplace
func TestReplace(t *testing.T) {
	source := `
package p
type S struct {
	Before string ` + "`" + `replaced:"false"` + "`" + `
}
func (s *S) String() string {
	return ` + "`" + `replaced:"false"` + "`" + `
}

func Hello() string {
	return ` + "`" + `replaced:"false"` + "`" + `
}
`
	type C struct {
		source  string
		source2 string
		name    string
		msg     string
		hasErr  bool
	}

	candidates := []C{
		{
			source: source,
			msg:    "replace struct",
			name:   "S",
			source2: `
package p
type S struct {
	After string ` + "`" + `replaced:"true"` + "`" + `
}`,
		},
		{
			source: source,
			msg:    "replace function",
			name:   "Hello",
			source2: `
package p
func Hello() string {
	return ` + "`" + `replaced:"true"` + "`" + `
}`,
		},
		{
			source: source,
			msg:    "replace method",
			name:   "S.String",
			source2: `
package p
func (s *S) String() string {
	return ` + "`" + `replaced:"true"` + "`" + `
}`,
		},
	}

	for _, c := range candidates {
		c := c
		t.Run(c.msg, func(t *testing.T) {
			pf := NewPatchwork().MustParseFile("f0", source)
			pf2 := pf.MustParseFile("f1", c.source2)

			t.Logf("input (%s)\n%s\n", c.name, c.source)
			t.Logf("replace (%s)\n%s\n", c.name, c.source2)

			ok, err := pf.Replace(pf2.Lookup(c.name))

			if c.hasErr {
				t.Logf("should error %s", err)
				if err == nil {
					t.Fatal("error is expected, but no error")
				}
				return
			}

			var b bytes.Buffer
			if err := pf.FprintCode(&b); err != nil {
				t.Fatal(err)
			}
			t.Logf("output\n%s\n", b.String())

			if err != nil {
				t.Fatal(err)
			}

			if !ok {
				t.Fatal("must replaceed")
			}

			// tentative assertion
			if !strings.Contains(b.String(), `replaced:"true"`) {
				t.Fatalf("cannot replaced (%q)", c.name)
			}
		})
	}
}
