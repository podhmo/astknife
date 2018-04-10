package patchwork5

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	candidates := []struct {
		msg            string
		source         string
		expectedShapes []string
	}{
		{
			msg:            "empty",
			source:         "package p",
			expectedShapes: []string{},
		},
		{
			msg: "struct-only",
			source: `
package p

// S0 :
type S0 struct {}

// S1 :
type S1 struct {}
`,
			expectedShapes: []string{
				`<declhead>`,
				`<specref name=S0>`,
				`<decltail>`,

				`<declhead>`,
				`<specref name=S1>`,
				`<decltail>`,
			},
		},
		{
			msg: "struct-paren",
			source: `
package p

type (
	// S0 :
	S0 struct{}
	// S1 :
	S1 struct{}
)
`,
			expectedShapes: []string{
				`<declhead>`,
				`<specref name=S0>`,
				`<specref name=S1>`,
				`<decltail>`,
			},
		},
		{
			msg: "func",
			source: `
package p

// Hello :
func Hello(){}

// Bye :
func Bye(){}
`,
			expectedShapes: []string{
				`<declref name="Hello">`,
				`<declref name="Bye">`,
			},
		},
		{
			msg: "func-with-comment",
			source: `
// header

package p

// middle

// Hello :
func Hello(){}

// middle2

// Bye :
func Bye(){}

// footer
`,
			expectedShapes: []string{
				`<commentref>`,
				`<declref name="Hello">`,
				`<commentref>`,
				`<declref name="Bye">`,
				`<commentref>`,
			},
		},
		{
			msg: "import",
			source: `
package p

import "fmt"
`,
			expectedShapes: []string{
				`<declhead>`,
				`<specref name=>`, // import
				`<decltail>`,
			},
		},
		{
			msg: "import-paren",
			source: `
package p

import (
	"fmt"
	"io"
)
`,
			expectedShapes: []string{
				`<declhead>`,
				`<specref name=>`, // import
				`<specref name=>`, // import
				`<decltail>`,
			},
		},
	}
	for _, c := range candidates {
		c := c
		t.Run(c.msg, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "", c.source, parser.ParseComments)

			require.NoError(t, err)
			file := parseASTFile(f)
			shapes := make([]string, len(file.Regions))
			for i := range file.Regions {
				shapes[i] = file.Regions[i].Ref.String()
			}
			assert.Exactly(t, c.expectedShapes, shapes)
		})
	}
}
