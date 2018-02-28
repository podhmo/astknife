package lookup

import (
	"go/ast"
	"strings"

	"github.com/pkg/errors"
)

// Lookup :
func Lookup(name string, files ...*ast.File) (*Result, error) {
	if strings.Contains(name, ".") {
		obAndMethod := strings.SplitN(name, ".", 2)
		return Method(obAndMethod[0], obAndMethod[1], files...)
	}
	return Toplevel(name, files...)
}

// Toplevel :
func Toplevel(name string, files ...*ast.File) (*Result, error) {
	for _, f := range files {
		ob := f.Scope.Lookup(name)
		if ob != nil {
			r := &Result{
				Name:   name,
				Object: ob,
				Type:   TypeToplevel,
				File:   f,
			}
			if ob.Type == ast.Fun {
				r.FuncDecl = ob.Decl.(*ast.FuncDecl)
			}
			return r, nil
		}
	}
	return nil, errors.Wrap(ErrNotFound, "toplevel")
}

var (
	// ErrNotFound :
	ErrNotFound = errors.New("not found")
)

// Method :
func Method(obname, name string, files ...*ast.File) (*Result, error) {
	obr, err := Toplevel(obname, files...)
	if err != nil {
		return nil, errors.WithMessage(err, "method recv")
	}

	ob := obr.Object

	for _, f := range files {
		for _, decl := range f.Decls {
			if t, ok := decl.(*ast.FuncDecl); ok {
				if !IsMethod(t) {
					continue
				}

				if t.Name.Name != name {
					continue
				}

				if !IsSameTypeOrPointer(ob, t.Recv.List[0].Type) {
					continue
				}

				return &Result{
					Name:     name,
					FuncDecl: t,
					Object:   ob,
					Type:     TypeMethod,
					File:     f,
				}, nil
			}
		}
	}
	return nil, errors.Wrap(ErrNotFound, "method func")
}

// AllMethods :
func AllMethods(obname string, files ...*ast.File) []*Result {
	rob, err := Toplevel(obname, files...)
	if err != nil {
		return nil
	}

	ob := rob.Object
	var r []*Result

	for _, f := range files {
		for _, decl := range f.Decls {
			if t, ok := decl.(*ast.FuncDecl); ok {
				if !IsMethod(t) {
					continue
				}

				if !IsSameTypeOrPointer(ob, t.Recv.List[0].Type) {
					continue
				}

				r = append(r, &Result{
					Name:     t.Name.Name,
					FuncDecl: t,
					Object:   ob,
					Type:     TypeMethod,
					File:     f,
				})
			}
		}
	}
	return r
}
