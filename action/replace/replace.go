package replace

import (
	"go/ast"
	"go/token"

	"github.com/pkg/errors"
	"github.com/podhmo/astknife/lookup"
)

// todo: all object types support
// // The list of possible Object kinds.
// const (
// 	Bad ObjKind = iota // for error handling
// 	Pkg                // package
// 	Con                // constant
// 	Typ                // type
// 	Var                // variable
// 	Fun                // function or method
// 	Lbl                // label
// )

/*
MEMO: for printer output correctly, surround with comments that replacement file.

e.g.

The situation, replacing f0.F to f1.F.

f0.go

```
// F : f0
func F() int {
	return 0
}
```

f1.go

```
// F : f1
func F() {
	return 1
}
```

After replace method called.

```
// F : f1 (comment of replacement)
func F() {
	return 1
} // (comment of replacement)
```

The reason of this, go/printer's printing method depends on physical position(token.Pos), and decision the timing of writing comment, is `current.Pos() > comment.Pos()`
*/

// Toplevel :
func Toplevel(fset *token.FileSet, dst *ast.File, dstOb *ast.Object, replacement *ast.Object) (ok bool, err error) {
	if replacement == nil {
		return
	}

	switch replacement.Kind {
	case ast.Typ:
		dstSpec, can := dstOb.Decl.(ast.Spec)
		if !can {
			err = errors.Errorf("invalid object type %s dst (kind=%q)", replacement.Type, replacement.Kind) // xxx
			return
		}
		t, can := replacement.Decl.(ast.Spec)
		if !can {
			err = errors.Errorf("invalid object type %s replacement (kind=%q)", replacement.Type, replacement.Kind) // xxx
			return
		}
		return Spec(fset, dst, dstSpec, t)
	case ast.Fun:
		dstDecl, can := dstOb.Decl.(*ast.FuncDecl)
		if !can {
			err = errors.Errorf("invalid object type %s dst (kind=%q)", replacement.Type, replacement.Kind) // xxx
			return
		}
		t, can := replacement.Decl.(*ast.FuncDecl)
		if !can {
			err = errors.Errorf("invalid object type %s replacement (kind=%q)", t.Type, replacement.Kind) // xxx
			return
		}
		return Function(fset, dst, dstDecl, t)
	default:
		err = errors.Errorf("unsupported object type %s (kind=%q)", replacement.Type, replacement.Kind)
		return
	}
}

// Spec :
func Spec(fset *token.FileSet, dst *ast.File, dstSpec ast.Spec, replacement ast.Spec) (ok bool, err error) {
	ast.Inspect(dst, func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.GenDecl:
			newspec := make([]ast.Spec, len(t.Specs))
			for i, spec := range t.Specs {
				if spec == dstSpec {
					ok = true
					newspec[i] = replacement
				} else {
					newspec[i] = spec
				}
			}
			t.Specs = newspec
			return false
			// case *ast.FuncDecl:
			// case *ast.BadDecl:
		}
		return true
	})
	return
}

// Function :
func Function(fset *token.FileSet, dst *ast.File, dstDecl *ast.FuncDecl, replacement *ast.FuncDecl) (ok bool, err error) {
	if replacement == nil {
		return
	}
	for i, decl := range dst.Decls {
		if decl == dstDecl {
			dst.Decls[i] = replacement
			ok = true
			return
		}
	}
	return
}

// Method :
func Method(fset *token.FileSet, dst *ast.File, ob *ast.Object, dstDecl *ast.FuncDecl, replacement *ast.FuncDecl) (ok bool, err error) {
	if replacement == nil {
		return
	}
	for i, decl := range dst.Decls {
		if decl, can := decl.(*ast.FuncDecl); can {
			if lookup.IsMethod(decl) && decl == dstDecl {
				dst.Decls[i] = replacement
				ok = true
			}
		}
	}
	return
}
