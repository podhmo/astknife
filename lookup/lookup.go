package lookup

import (
	"go/ast"
	"strings"
)

// Lookup :
type Lookup struct {
	lookup func(name string) *ast.Object
	Files  []*ast.File
}

// New :
func New(files ...*ast.File) *Lookup {
	k := &Lookup{Files: files}
	k.lookup = func(name string) *ast.Object {
		for _, f := range k.Files {
			if ob := f.Scope.Lookup(name); ob != nil {
				return ob
			}
		}
		return nil
	}
	return k
}

// With :
func (k *Lookup) With(file *ast.File) *Lookup {
	files := []*ast.File{file}
	files = append(files, k.Files...)
	return New(files...)
}

// Lookup :
func (k *Lookup) Lookup(name string) *Result {
	if strings.Contains(name, ".") {
		obAndMethod := strings.SplitN(name, ".", 2)
		return k.Method(obAndMethod[0], obAndMethod[1])
	}
	return k.Toplevel(name)
}

// Toplevel :
func (k *Lookup) Toplevel(name string) *Result {
	raw := k.lookup(name)
	if raw == nil {
		return nil
	}
	return &Result{
		Type:   TypeToplevel,
		Object: raw,
	}
}

// AllMethods :
func (k *Lookup) AllMethods(obname string) []*Result {
	ob := k.lookup(obname)
	if ob == nil {
		return nil
	}

	var r []*Result
	for _, f := range k.Files {
		for _, decl := range f.Decls {
			if decl, ok := decl.(*ast.FuncDecl); ok {
				if IsMethod(decl) && IsSameTypeOrPointer(ob, decl.Recv.List[0].Type) {
					r = append(r, &Result{
						Type:     TypeMethod,
						FuncDecl: decl,
					})
				}
			}
		}
	}
	return r
}

// Method :
func (k *Lookup) Method(obname string, name string) *Result {
	ob := k.lookup(obname)
	if ob == nil {
		return nil
	}
	return k.MethodByObject(ob, name)
}

// MethodByObject :
func (k *Lookup) MethodByObject(ob *ast.Object, name string) *Result {
	for _, f := range k.Files {
		for _, decl := range f.Decls {
			if decl, ok := decl.(*ast.FuncDecl); ok {
				if IsMethod(decl) && IsSameTypeOrPointer(ob, decl.Recv.List[0].Type) {
					if decl.Name.Name == name {
						return &Result{
							Type:     TypeMethod,
							Object:   ob,
							FuncDecl: decl,
						}
					}
				}
			}
		}
	}
	return nil
}
