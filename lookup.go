package astknife

import (
	"go/ast"
)

// LookupType :
type LookupType string

const (
	// LookupTypeToplevel : toplevel (e.g. toplevel function, struct definition)
	LookupTypeToplevel = LookupType("toplevel")
	// LookupTypeMethod : method (e.g. method function, struct definition)
	LookupTypeMethod = LookupType("method")
)

// LookupResult :
type LookupResult struct {
	Type     LookupType
	FuncDecl *ast.FuncDecl
	Object   *ast.Object
}

// Name :
func (r *LookupResult) Name() string {
	switch r.Type {
	case LookupTypeToplevel:
		return r.Object.Name
	case LookupTypeMethod:
		return r.FuncDecl.Name.Name
	}
	return "<nil>"
}

// LookupToplevel :
func LookupToplevel(file *ast.File, name string) *LookupResult {
	raw := file.Scope.Lookup(name)
	if raw == nil {
		return nil
	}
	return &LookupResult{
		Type:   LookupTypeToplevel,
		Object: raw,
	}
}

// LookupAllMethods :
func LookupAllMethods(f *ast.File, obname string) []*LookupResult {
	ob := f.Scope.Lookup(obname)
	if ob == nil {
		return nil
	}

	var r []*LookupResult
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok {
			if IsMethod(decl) && IsSameTypeOrPointer(ob, decl.Recv.List[0].Type) {
				r = append(r, &LookupResult{
					Type:     LookupTypeMethod,
					FuncDecl: decl,
				})
			}
		}
	}
	return r
}

// LookupMethod :
func LookupMethod(f *ast.File, obname string, name string) *LookupResult {
	ob := f.Scope.Lookup(obname)
	if ob == nil {
		return nil
	}

	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok {
			if IsMethod(decl) && IsSameTypeOrPointer(ob, decl.Recv.List[0].Type) {
				if decl.Name.Name == name {
					return &LookupResult{
						Type:     LookupTypeMethod,
						FuncDecl: decl,
					}
				}
			}
		}
	}
	return nil
}
