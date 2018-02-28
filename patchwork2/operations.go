package patchwork2

import (
	"go/ast"
	"go/token"

	"github.com/pkg/errors"
	"github.com/podhmo/astknife/patchwork2/lookup"
)

// Append :
func Append(fref *FileRef, r *lookup.Result) {
	fref.Decls = append(fref.Decls, newDeclFromLookup(r))
}

func newDeclFromLookup(r *lookup.Result) declRef {
	switch r.Object.Kind {
	case ast.Typ:
        return &GenDeclRef{
            Replacement: xx,
        }
		dst.Decls = append(dst.Decls, &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{ob.Decl.(ast.Spec)},
		})
		ok = true
	case ast.Fun:
		if decl, can := ob.Decl.(*ast.FuncDecl); can {
			return FunctionToFile(dst, decl)
		}
		err = errors.Errorf("unsupported object type %s (kind=%q)", ob.Type, ob.Kind)
		return
	default:
		err = errors.Errorf("unsupported object type %s (kind=%q)", ob.Type, ob.Kind)
		return
	}
}
