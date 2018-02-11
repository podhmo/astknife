package lookup

import (
	"go/ast"
)

// Type :
type Type string

const (
	// TypeToplevel : toplevel (e.g. toplevel function, struct definition)
	TypeToplevel = Type("toplevel")
	// TypeMethod : method (e.g. method function, struct definition)
	TypeMethod = Type("method")
)

// Result :
type Result struct {
	Type     Type
	FuncDecl *ast.FuncDecl
	Object   *ast.Object
}

// Name :
func (r *Result) Name() string {
	switch r.Type {
	case TypeToplevel:
		return r.Object.Name
	case TypeMethod:
		return r.FuncDecl.Name.Name
	}
	return "<nil>"
}

// Toplevel :
func Toplevel(file *ast.File, name string) *Result {
	raw := file.Scope.Lookup(name)
	if raw == nil {
		return nil
	}
	return &Result{
		Type:   TypeToplevel,
		Object: raw,
	}
}

// AllMethods :
func AllMethods(f *ast.File, obname string) []*Result {
	ob := f.Scope.Lookup(obname)
	if ob == nil {
		return nil
	}

	var r []*Result
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok {
			if isMethod(decl) && isSameTypeOrPointer(ob, decl.Recv.List[0].Type) {
				r = append(r, &Result{
					Type:     TypeMethod,
					FuncDecl: decl,
				})
			}
		}
	}
	return r
}

// Method :
func Method(f *ast.File, obname string, name string) *Result {
	ob := f.Scope.Lookup(obname)
	if ob == nil {
		return nil
	}
	return MethodFromObject(f, ob, name)
}

// MethodFromObject :
func MethodFromObject(f *ast.File, ob *ast.Object, name string) *Result {
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok {
			if isMethod(decl) && isSameTypeOrPointer(ob, decl.Recv.List[0].Type) {
				if decl.Name.Name == name {
					return &Result{
						Type:     TypeMethod,
						FuncDecl: decl,
					}
				}
			}
		}
	}
	return nil
}
