package patchwork2

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/podhmo/astknife/patchwork2/lookup"
)

var errNotfound = errors.New("not found")

// Replace :
func Replace(ref *Ref, pat *lookup.Result, rep *lookup.Result) error {
	if pat.Object.Kind != rep.Object.Kind {
		return &Unsupported{pat, fmt.Sprintf("conflict kind %s != %s", pat.Object.Kind, rep.Object.Kind)}
	}
	if rep == nil {
		return errNotfound
	}

	for _, f := range ref.Files {
		err := replaceFileRefFromLookup(f, pat, rep)
		if err == nil {
			return nil
		}
		if err == errNotfound {
			continue
		}
		return err
	}
	return &Unsupported{pat, "replace target not found"}
}

func replaceFileRefFromLookup(fref *FileRef, pat *lookup.Result, rep *lookup.Result) error {
	for _, f := range fref.Decls {
		err := replaceDeclRefFromLookup(f, pat, rep)
		if err == nil {
			return nil
		}
		if err == errNotfound {
			continue
		}
		return err
	}
	return errNotfound
}

func replaceDeclRefFromLookup(dref *DeclRef, pat *lookup.Result, rep *lookup.Result) error {
	switch pat.Object.Kind {
	case ast.Typ:
		spec, ok := pat.Object.Decl.(ast.Spec)
		if !ok {
			return &Unsupported{pat, "not spec"}
		}
		for _, sref := range dref.Specs {
			if sref.Original == spec {
				sref.Replacement = rep.Object.Decl.(ast.Spec)
				sref.File = rep.File
				sref.Result = rep
				return nil
			}
		}
		return errNotfound
	case ast.Fun:
		decl, ok := pat.Object.Decl.(ast.Decl)
		if !ok {
			return &Unsupported{pat, "not decl"}
		}
		if dref.Original == decl {
			dref.Replacement = rep.Object.Decl.(ast.Decl)
			dref.File = rep.File
			dref.Result = rep
			return nil
		}
		return errNotfound
	default:
		return &Unsupported{pat, "invalid kind"}
	}
}

// Append :
func Append(fref *FileRef, r *lookup.Result) error {
	if r == nil {
		return errNotfound
	}
	declref, err := newDeclRefFromLookup(r)
	if err != nil {
		return err
	}
	fref.Decls = append(fref.Decls, declref)
	return nil
}

func newDeclRefFromLookup(r *lookup.Result) (*DeclRef, error) {
	switch r.Object.Kind {
	case ast.Typ:
		spec, ok := r.Object.Decl.(ast.Spec)
		if !ok {
			return nil, &Unsupported{r, "not spec"}
		}
		var tok token.Token
		switch spec.(type) {
		case *ast.ImportSpec:
			tok = token.IMPORT
		case *ast.TypeSpec:
			tok = token.TYPE
		// case *ast.ValueSpec:
		// 	tok = token.CONST
		case *ast.ValueSpec:
			tok = token.VAR
		default:
			return nil, &Unsupported{r, "invalid spec"}
		}
		decl := &ast.GenDecl{
			Tok:    tok,
			Specs:  []ast.Spec{spec},
			TokPos: spec.Pos(), // xxx:
		}
		return &DeclRef{
			Replacement: decl,
			File:        r.File,
			Result:      r,
			Specs: []*SpecRef{&SpecRef{
				Replacement: spec,
				File:        r.File,
			}},
		}, nil
	case ast.Fun:
		return &DeclRef{
			Replacement: r.FuncDecl,
			Result:      r,
			File:        r.File,
		}, nil
	default:
		return nil, &Unsupported{r, "invalid kind"}
	}
}

// Unsupported :
type Unsupported struct {
	Result *lookup.Result
	Reason string
}

// Unsupported :
func (e *Unsupported) Error() string {
	return fmt.Sprintf("unsupported object %s (reason=%s, type=%s, kind=%q)",
		e.Result.String(), e.Reason, e.Result.Object.Type, e.Result.Object.Kind,
	)
}
