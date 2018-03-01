package lookup

import (
	"go/ast"
	"strings"
)

// Lookup :
func Lookup(name string, files ...*ast.File) *Result {
	if strings.Contains(name, ".") {
		obAndMethod := strings.SplitN(name, ".", 2)
		return Method(obAndMethod[0], obAndMethod[1], files...)
	}
	return Toplevel(name, files...)
}

// Toplevel :
func Toplevel(name string, files ...*ast.File) *Result {
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
			return r
		}
	}
	return nil
}

// Method :
func Method(obname, name string, files ...*ast.File) *Result {
	obr := Toplevel(obname, files...)
	if obr == nil {
		return nil
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
				}
			}
		}
	}
	return nil
}

// AllMethods :
func AllMethods(obname string, files ...*ast.File) []*Result {
	rob := Toplevel(obname, files...)
	if rob == nil {
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
