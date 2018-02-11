package patchwork

import (
	"go/ast"

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

// todo: comment support

// replaceToplevelToFile :
func replaceToplevelToFile(dst *ast.File, dstOb *ast.Object, ob *ast.Object) (ok bool, err error) {
	if ob == nil {
		return
	}

	switch ob.Kind {
	case ast.Typ:
		dstSpec, can := dstOb.Decl.(ast.Spec)
		if !can {
			err = errors.Errorf("invalid object type %s dst (kind=%q)", ob.Type, ob.Kind) // xxx
			return
		}
		replacement, can := ob.Decl.(ast.Spec)
		if !can {
			err = errors.Errorf("invalid object type %s replacement (kind=%q)", ob.Type, ob.Kind) // xxx
			return
		}
		return replaceSpecToFile(dst, dstSpec, replacement)
	case ast.Fun:
		dstDecl, can := dstOb.Decl.(*ast.FuncDecl)
		if !can {
			err = errors.Errorf("invalid object type %s dst (kind=%q)", ob.Type, ob.Kind) // xxx
			return
		}
		replacement, can := ob.Decl.(*ast.FuncDecl)
		if !can {
			err = errors.Errorf("invalid object type %s replacement (kind=%q)", ob.Type, ob.Kind) // xxx
			return
		}
		return replaceFunctionToFile(dst, dstDecl, replacement)
	default:
		err = errors.Errorf("unsupported object type %s (kind=%q)", ob.Type, ob.Kind)
		return
	}
}

// replaceSpecToFile :
func replaceSpecToFile(dst *ast.File, dstSpec ast.Spec, replacement ast.Spec) (ok bool, err error) {
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

// replaceFunctionToFile :
func replaceFunctionToFile(dst *ast.File, dstDecl *ast.FuncDecl, replacement *ast.FuncDecl) (ok bool, err error) {
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

// replaceMethodToFile :
func replaceMethodToFile(dst *ast.File, ob *ast.Object, dstDecl *ast.FuncDecl, replacement *ast.FuncDecl) (ok bool, err error) {
	if replacement == nil {
		return
	}
	for i, decl := range dst.Decls {
		if decl, can := decl.(*ast.FuncDecl); can {
			if lookup.IsMethod(decl) && lookup.IsSameTypeOrPointer(ob, decl.Recv.List[0].Type) {
				dst.Decls[i] = replacement
				ok = true
			}
		}
	}
	return
}
